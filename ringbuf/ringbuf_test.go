package ringbuf

import "fmt"

func Example() {
	var rb RingBuf

	buf := make([]byte, 1000)

	fmt.Println("writing 5 x 1000 bytes")
	fmt.Println(rb.Write(buf))
	fmt.Println(rb.Write(buf))
	fmt.Println(rb.Write(buf))
	fmt.Println(rb.Write(buf))
	fmt.Println(rb.Write(buf))

	fmt.Println("reading 1000 and then writing 1000 more")
	fmt.Println(rb.Read(buf))
	fmt.Println(rb.Write(buf))

	fmt.Println("reading 5 x 1000 bytes")
	fmt.Println(rb.Read(buf))
	fmt.Println(rb.Read(buf))
	fmt.Println(rb.Read(buf))
	fmt.Println(rb.Read(buf))
	fmt.Println(rb.Read(buf))

	fmt.Println("reading one additional time")
	fmt.Println(rb.Read(buf))

	// Output:
	//
	// writing 5 x 1000 bytes
	// 1000 <nil>
	// 1000 <nil>
	// 1000 <nil>
	// 1000 <nil>
	// 96 <nil>
	// reading 1000 and then writing 1000 more
	// 1000 <nil>
	// 1000 <nil>
	// reading 5 x 1000 bytes
	// 1000 <nil>
	// 1000 <nil>
	// 1000 <nil>
	// 1000 <nil>
	// 96 <nil>
	// reading one additional time
	// 0 EOF
}
