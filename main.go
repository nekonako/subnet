package main

import (
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

func main() {
	var ip = os.Args[1]

	address, prefix, totalAddr, networkAddr, hostAddr, broadcastAddr := calc(ip)

	fmt.Println("")
	fmt.Println("ip address           :", address)
	fmt.Println("prefix               :", prefix)
	fmt.Println("total address        :", totalAddr)
	fmt.Println("network address      :")
	for _, value := range networkAddr {
		fmt.Println("                      ", value)
	}
	fmt.Println("host address         :")
	for _, value := range hostAddr {
		fmt.Println("                      ", value)
	}
	fmt.Println("broadcast address    :")
	for _, value := range broadcastAddr {
		fmt.Println("                      ", value)
	}
}

func splitAddress(ip string) (address string, prefix int, class string) {
	split := strings.Split(ip, "/")
	if len(split) <= 1 {
		panic(fmt.Errorf("Error : prefix not found"))
	}
	prefix, _ = strconv.Atoi(split[1])
	if prefix >= 24 && prefix <= 30 {
		class = "c"
	} else if prefix >= 16 && prefix <= 30 {
		class = "b"
	} else if prefix >= 8 && prefix <= 30 {
		class = "a"
	} else {
		fmt.Println("Prefix must range of 8 - 30")
		panic("ip address class not found")
	}
	return split[0], prefix, class
}

func splitIp(ip string, class string) (subnet string) {
	addr := strings.Split(ip, ".")
	if class == "c" {
		subnet = addr[0] + "." + addr[1] + "." + addr[2]
	} else if class == "b" {
		subnet = addr[0] + "." + addr[1]
	} else if class == "a" {
		subnet = addr[0]
	}
	return subnet
}

func subnetBlock(prefix int) int {
	var bitList = [8]int{0, 128, 64, 32, 16, 8, 4, 2}
	var number int
	if prefix >= 24 {
		bit := prefix - 24
		for i := 0; i <= bit; i++ {
			number += bitList[i]
		}
	} else if prefix >= 16 {
		bit := prefix - 16
		for i := 0; i <= bit; i++ {
			number += bitList[i]
		}
	} else if prefix >= 8 {
		bit := prefix - 8
		for i := 0; i <= bit; i++ {
			number += bitList[i]
		}
	} else {
		panic(fmt.Errorf("Prefix is not valid"))
	}
	return 256 - number
}

func getNetworkAddress(subnetBlock int) (networkAddr []int) {
	for i := 0; i < 255; i++ {
		networkAddr = append(networkAddr, i)
		i += int(subnetBlock - 1)
	}
	return networkAddr
}

func getHostAddress(subnetBlock int, networkAddr []int) (hostAddr [][]int) {
	for _, val := range networkAddr {
		var blok []int
		for i := val + 1; i < val+subnetBlock; i++ {
			blok = append(blok, i)
		}
		hostAddr = append(hostAddr, blok)
	}
	return hostAddr
}

func getBroadcasAddress(subnetBlock int, networkAddr []int) (broadcastAddr []int) {
	for _, val := range networkAddr {
		if val != 0 {
			broadcastAddr = append(broadcastAddr, val-1)
		}
	}
	broadcastAddr = append(broadcastAddr, 255)
	return broadcastAddr
}

func calc(ip string) (address string, prefix int, totalAddr float64, networkAddr []string, hostAddr []string, broadcastAddr []string) {

	address, prefix, class := splitAddress(ip)
	totalAddr = math.Pow(2, float64(32-prefix))
	subnetBlock := subnetBlock(prefix)
	subnet := splitIp(address, class)
	network := getNetworkAddress(subnetBlock)
	host := getHostAddress(subnetBlock, network)
	broadcast := getBroadcasAddress(subnetBlock, network)

	for _, value := range network {
		if class == "c" {
			networkAddr = append(networkAddr, subnet+"."+strconv.Itoa(value))
		} else if class == "b" {
			networkAddr = append(networkAddr, subnet+"."+strconv.Itoa(value)+".0")
		} else if class == "a" {
			networkAddr = append(networkAddr, subnet+"."+strconv.Itoa(value)+".0.0")
		}
	}

	for _, value := range host {
		if len(value) >= 2 {
			if class == "c" {
				hostAddr = append(hostAddr, subnet+"."+strconv.Itoa(value[0])+" - "+subnet+"."+strconv.Itoa(value[len(value)-2]))
			} else if class == "b" {
				hostAddr = append(hostAddr, subnet+"."+strconv.Itoa(value[0])+".0"+" - "+subnet+"."+strconv.Itoa(value[len(value)-2])+".255")
			} else if class == "a" {
				hostAddr = append(hostAddr, subnet+"."+strconv.Itoa(value[0])+".0.0"+" - "+subnet+"."+strconv.Itoa(value[len(value)-2])+".255.255")
			}
		} else {
			if class == "c" {
				hostAddr = append(hostAddr, subnet+"."+strconv.Itoa(value[0]))
			} else if class == "b" {
				hostAddr = append(hostAddr, subnet+"."+strconv.Itoa(value[0])+".0"+" - "+subnet+"."+strconv.Itoa(value[0])+".255")
			} else if class == "a" {
				hostAddr = append(hostAddr, subnet+"."+strconv.Itoa(value[0])+".0.0"+" - "+subnet+"."+strconv.Itoa(value[0])+".255.255")
			}
		}
	}

	for _, value := range broadcast {
		if class == "c" {
			broadcastAddr = append(broadcastAddr, subnet+"."+strconv.Itoa(value))
		}
		if class == "b" {
			broadcastAddr = append(broadcastAddr, subnet+"."+strconv.Itoa(value)+".255")
		}
		if class == "a" {
			broadcastAddr = append(broadcastAddr, subnet+"."+strconv.Itoa(value)+".255.255")
		}
	}

	return address, prefix, totalAddr, networkAddr, hostAddr, broadcastAddr

}
