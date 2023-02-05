package cachekeys

const (
	//VehicleRecordListCacheKey defines the cache key for storing vehicle record list
	VehicleRecordListCacheKey string = "vehicle-list"
	//VehicleRecordCacheKey defines the cache key for storing vehicle record. %s stands for the vehicle stock ID
	VehicleRecordCacheKey = "vehicle:%s"
	//BasicRecordCacheChildKey defines the hash sub key for storing basic record
	BasicRecordCacheChildKey = "basic-record"
	//DetailedRecordCacheChildKey defines the hash sub key for storing detailed record
	DetailedRecordCacheChildKey = "detailed-record"
	//VehicleReviewsCacheKey defines the cache key for storing reviews list. %s stands for the vehicle stock ID
	VehicleReviewsCacheKey = "vehicle:%s:reviews"
)
