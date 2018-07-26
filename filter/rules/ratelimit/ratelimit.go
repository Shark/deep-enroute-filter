package ratelimit

import (
  "crypto/sha256"
  "fmt"
  "time"

  "gitlab.hpi.de/felix.seidel/iotsec-enroute-filtering/filter/types"

  "github.com/Clever/leakybucket"
  "github.com/Clever/leakybucket/memory"
)

type rateLimitStateEvent struct {
  payload []bucketInfo
}

func (e *rateLimitStateEvent) Type() string {
  return "RateLimitStateEvent"
}

func (e *rateLimitStateEvent) Payload() interface{} {
  return e.payload
}

type bucketInfo struct {
  SrcIP string
  DstIP string
  BucketCapacity uint
  BucketRemaining uint
  BucketReset time.Time
}

type RateLimit struct {
  events chan types.Event
  buckets *memory.Storage
  bucketKeys map[string]bucketInfo
}

func (r RateLimit) Name() string {
  return "RateLimit"
}

func (r RateLimit) Process(message *types.COAPMessage) types.RuleProcessingResult {
  key := messageInfo(message)
  duration, _ := time.ParseDuration("10s")
  bucket, err := r.buckets.Create(key, 10, duration)

  if err != nil {
    message := fmt.Sprintf("Error accessing bucket: %v", err)
    return types.RuleProcessingResult{false, r, &message}
  }

  state, err := bucket.Add(1)

  r.bucketKeys[key] = bucketInfo{
    message.Metadata.SrcIP,
    message.Metadata.DstIP,
    state.Capacity,
    state.Remaining,
    state.Reset,
  }

  if err == nil {
    return types.RuleProcessingResult{
      true,
      r,
      nil,
    }
  } else if err == leakybucket.ErrorFull {
    message := fmt.Sprintf("Bucket is full: %d/%d", state.Remaining, state.Capacity)
    return types.RuleProcessingResult{
      false,
      r,
      &message,
    }
  } else {
    message := fmt.Sprintf("Error adding to bucket: %v", err)
    return types.RuleProcessingResult{
      false,
      r,
      &message,
    }
  }
}

func (r RateLimit) publishState() {
  bucketInfos := make([]bucketInfo, 0)

  for _, bucketInfo := range r.bucketKeys {
    info := bucketInfo
    if bucketInfo.BucketReset.Before(time.Now()) {
      info.BucketRemaining = info.BucketCapacity
    }
    bucketInfos = append(bucketInfos, info)
  }

  r.events <- &rateLimitStateEvent{bucketInfos}
}

func NewRateLimit(events chan types.Event) *RateLimit {
  rateLimit := &RateLimit{
    events,
    memory.New(),
    make(map[string]bucketInfo),
  }

  go func() {
    ticker := time.NewTicker(1 * time.Second)
    for _ = range ticker.C {
      rateLimit.publishState()
    }
  }()

  return rateLimit
}

func messageInfo(message *types.COAPMessage) string {
  stringKey := fmt.Sprintf("%s%s", message.Metadata.SrcIP, message.Metadata.DstIP)
  byteSum := sha256.Sum256([]byte(stringKey))
  return fmt.Sprintf("%x", byteSum)
}
