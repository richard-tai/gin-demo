package priorityqueue

import (
	"demo/util"

	"github.com/redis/go-redis/v9"
)

type PriorityQueue struct {
	Name string
	Rcc  *redis.ClusterClient
}

func New(name string, rcc *redis.ClusterClient) *PriorityQueue {
	return &PriorityQueue{
		Name: name,
		Rcc:  rcc,
	}
}

const MillisecondsBase int64 = 1e14

// in ms
func GetTimeScore(timestamp int64) int64 {
	return MillisecondsBase - 1 - timestamp
}

func GetNowScore() int64 {
	return GetTimeScore(util.GetNowMilli())
}

// priority support [0,9] only
func GetScore(priority int64) int64 {
	return (priority%10)*MillisecondsBase + GetNowScore()
}

func GetPriority(score int64) int64 {
	return score / MillisecondsBase
}

func GetTime(score int64) int64 {
	return MillisecondsBase - 1 - score%MillisecondsBase
}

func (pq *PriorityQueue) EnqueueIfNX(keys []string, score int64, value string, useHigherScore bool) error {
	ss := `
		local key = KEYS[1]
		local score = ARGV[1]
		local member = ARGV[2]
		for i = 2, #KEYS do
			local one_key = KEYS[i]
			if redis.call('exists', one_key) ~= 0 and redis.call('zrank', one_key, member) then
				return redis.call('zrevrank', one_key, member)
			end
			
			if redis.call('exists', key) == 0 or not redis.call('zrank', key, member) then
				redis.call('zadd', key, score, member)
			else
				if ARGV[3] == "useHigherScore" then
					local old_score = redis.call('zscore', key, member)
					if score > old_score then
						redis.call('zadd', key, score, member)
					end
				end
			end
			
			return redis.call('zrevrank', key, member)
		end
	`
	useHigherScoreStr := ""
	if useHigherScore {
		useHigherScoreStr = "useHigherScore"
	}
	_, err := redis.NewScript(ss).Run(pq.Rcc, append([]string{pq.Name}, keys...), score, value, useHigherScoreStr).Result()
	return err
}

func (pq *PriorityQueue) Enqueue(score int64, value string) error {
	return pq.EnqueueIfNX([]string{pq.Name}, score, value, true)
}
