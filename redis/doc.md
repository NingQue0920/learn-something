1. 尽可能使用Hash类型，Redis会对Hash类型数据选取合理的编码格式，从而提高内存利用率和查询性能。
   - 数据量小，Hash使用ziplist/listpack进行编码，压缩内存，同时保证线性查询的时间复杂度。
   - 数据量大，Hash使用hashtable进行编码，保证了K-V映射关系，查询时间复杂度为O(1)

