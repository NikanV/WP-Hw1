package __

import (
	"io/ioutil"
	"math"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))
var primes []int = nil

func RandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyz" + "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

/////////////////////////////////////////////////////////////////////////random num

func RandomNumber(max int) int64 {
	return int64(seededRand.Intn(max))
}

func readFile(fname string) (nums []int, err error) {
	b, err := ioutil.ReadFile(fname)
	if err != nil {
		return nil, err
	}

	lines := strings.Fields(string(b))
	nums = make([]int, 0, len(lines))

	for _, l := range lines {
		if len(l) == 0 {
			continue
		}
		n, err := strconv.Atoi(l)
		if err != nil {
			return nil, err
		}
		nums = append(nums, n)
	}

	return nums, nil
}

func Random_Prime() int {
	if primes == nil {
		primes, _ = readFile("../tools/prime.txt")
	}
	return primes[seededRand.Intn(len(primes))]
}

////////////////////////////////////////////////////////////////////////// primitive root

func isPrime(n int) bool {
	// Corner cases
	if n <= 1 {
		return false
	}
	if n <= 3 {
		return true
	}

	// This is checked so that we can skip
	// middle five numbers in below loop
	if n%2 == 0 || n%3 == 0 {
		return false
	}
	for i := 5; i*i <= n; i = i + 6 {
		if n%i == 0 || n%(i+2) == 0 {
			return false
		}
	}

	return true
}

func power(x int, y uint, p int) int {
	res := 1 // Initialize result

	x = x % p // Update x if it is more than or
	// equal to p

	for { ///////////maybe problematic
		// If y is odd, multiply x with result
		if y&1 == 1 {
			res = (res * x) % p
		}
		// y must be even now
		y = y >> 1 // y = y/2
		x = (x * x) % p
		if y <= 0 {
			break
		}
	}
	return res
}

func findPrimefactors(s map[int]struct{}, n int) {
	// Print the number of 2s that divide n
	for {
		s[2] = struct{}{}
		n = n / 2
		if n%2 == 0 {
			break
		}
	}

	// n must be odd at this point. So we can skip
	// one element (Note i = i +2)
	for i := 3; i <= int(math.Sqrt(float64(n))); i = i + 2 {
		// While i divides n, print i and divide n
		for {
			s[i] = struct{}{}
			n = n / i
			if n%i == 0 {
				break
			}
		}
	}

	// This condition is to handle the case when
	// n is a prime number greater than 2
	if n > 2 {
		s[n] = struct{}{}
	}
}

func FindPrimitive(n int) int {
	s := map[int]struct{}{}

	// Check if n is prime or not
	if isPrime(n) == false {
		return -1
	}

	// Find value of Euler Totient function of n
	// Since n is a prime number, the value of Euler
	// Totient function is n-1 as there are n-1
	// relatively prime numbers.
	phi := n - 1

	// Find prime factors of phi and store in a set
	findPrimefactors(s, phi)

	// Check for every number from 2 to phi
	for r := 2; r <= phi; r++ {
		// Iterate through all prime factors of phi.
		// and check if we found a power with value 1
		flag := false
		for it := range s {

			// Check if r^((phi)/primefactors) mod n
			// is 1 or not
			if power(r, uint(phi/(it)), n) == 1 {
				flag = true
				break
			}
		}

		// If there was no power with value 1.
		if flag == false {
			return r
		}
	}

	// If no primitive root found
	return -1
}
