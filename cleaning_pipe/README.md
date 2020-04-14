### cleaning pipe

wraps in io.Reader and cleans unwanted bytes from the stream

| benchmark        | ops           | ns/op  |
| ------------- |-------------| -----|
| BenchmarkCleanBodyPipe_Read100-4      |                   5033256      |         229 ns/op  |           544 B/op          2 allocs/op|
|BenchmarkCleanBodyPipe_Read1000-4                 |       4947784         |      229 ns/op       |      544 B/op     |     2 allocs/op|
|BenchmarkCleanBodyPipe_Read100000-4               |       5020784        |       238 ns/op       |      547 B/op     |     2 allocs/op|
|BenchmarkBaseLineReadAllAndReplace_Read100-4      |       5268264        |       226 ns/op       |      512 B/op     |     1 allocs/op|
|BenchmarkBaseLineReadAllAndReplace_Read1000-4     |       5260072        |       227 ns/op       |      512 B/op     |     1 allocs/op|
|BenchmarkBaseLineReadAllAndReplace_Read100000-4   |       5046278        |       238 ns/op       |      515 B/op     |     1 allocs/op|


