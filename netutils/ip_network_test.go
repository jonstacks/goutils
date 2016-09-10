package netutils

import (
  "fmt"
  "net"
  "testing"
)

var ip, ip_net, net_err = net.ParseCIDR("192.168.1.2/24")


func TestIPNetwork(t *testing.T) {
  // Create a new channel to read from
  c := IPNetwork(ip_net)
  a := <-c
  if !a.Equal(net.ParseIP("192.168.1.0")) {
    t.Error("Does not match")
  }
}

func BenchmarkIPNetwork(b *testing.B) {
  for i := 0; i < b.N; i++ {
    for c := range IPNetwork(ip_net) {
      fmt.Println(c)
    }
  }
}

func ExampleIPNetwork() {
  c := IPNetwork(ip_net)
  fmt.Println(<- c)
  fmt.Println(<- c)
  fmt.Println(<- c)
  fmt.Println(<- c)
  // Output:
  // 192.168.1.0
  // 192.168.1.1
  // 192.168.1.2
  // 192.168.1.3
}
