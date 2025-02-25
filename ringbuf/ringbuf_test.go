package ringbuf

import "fmt"

func Example() {
	var rb RingBuf

	buf := make([]byte, 1000)

	fmt.Println(rb.Write(buf))
	fmt.Println(rb.Write(buf))
	fmt.Println(rb.Write(buf))
	fmt.Println(rb.Write(buf))
	fmt.Println(rb.Write(buf))

	fmt.Println(rb.Read(buf))
	fmt.Println(rb.Read(buf))
	fmt.Println(rb.Read(buf))
	fmt.Println(rb.Read(buf))
	fmt.Println(rb.Read(buf))

	fmt.Println(rb.Read(buf))

	// Output:
	//
	// 1000 <nil>
	// 1000 <nil>
	// 1000 <nil>
	// 1000 <nil>
	// 96 <nil>
	// 1000 <nil>
	// 1000 <nil>
	// 1000 <nil>
	// 1000 <nil>
	// 96 <nil>
	// 0 EOF
}
