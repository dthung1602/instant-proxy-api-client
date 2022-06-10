<!-- README template from https://github.com/dthung1602/instant-proxy-api-client -->


[![Contributors][contributors-shield]][contributors-url]
[![Forks][forks-shield]][forks-url]
[![Stargazers][stars-shield]][stars-url]
[![Issues][issues-shield]][issues-url]
[![MIT License][license-shield]][license-url]

<!-- PROJECT LOGO -->
<br />
<p align="center">
  <a href="https://github.com/dthung1602/instant-proxy-api-client">
    <img src="./resources/instantproxies.png" width="500">
  </a>

<h3 align="center">InstantProxy API Client</h3>

<p align="center">
   An <b>unofficial</b>  go client for Instant Proxy API
</p>


## Installation

```bash
go get github.com/dthung1602/instant-proxy-api-client
```

## Usage

### Make API calls

```go
import (
    "net"
    ipac "github.com/dthung1602/instant-proxy-api-client"
)

// create client
endpoint := ""
client := ipac.NewClient("username", "password", endpoint)

// If endpoint is empty sting, the client make http call
// to the real instant proxy server at https://admin.instantproxies.com
// For testing, start a fake server and pass endpoint = "http://localhost:3000"
// See "Fake server" section

// You don't need to explicitly call Authenticate() before making API calls
// The client will do that automatically
// However, you can call Authenticate to check if the logins is correct 
err := client.Authenticate()

// Read API
proxies, err := client.GetProxies()
ips, err := client.GetAuthorizedIPs()

// Write API
err := client.AddAuthorizedIP(net.ParseIP("127.0.0.1"))
err := client.AddAuthorizedIPs([]net.IP{
    net.ParseIP("132.4.4.6"),
    net.ParseIP("56.77.3.2"),
})

err = client.RemoveAuthorizedIP(net.ParseIP("127.0.0.1"))
err = client.RemoveAuthorizedIPs([]net.IP{
    net.ParseIP("132.4.4.6"),
    net.ParseIP("56.77.3.2"),
})

err = client.SetAuthorizedIPs([]net.IP{
    net.ParseIP("8.8.8.8"),
    net.ParseIP("4.4.4.4"),
})

// Other
myIP, err := client.GetOwnedPublicIP()

proxy, err := ipac.MakeProxy("123.21.1.1:3000")

testProxies, err := ipac.MakeProxies([]string{
    "154.38.148.102:8800",
    "154.37.249.170:8800",
    "87.101.80.161:8800",
})
result := client.TestProxies(testProxies)

```

### Fake server

This module also includes a dead simple replica of InstantProxy server.
Login credential:
- Username: `username`
- Password: `password`

```go
import (
    ipac "github.com/dthung1602/instant-proxy-api-client"
)

port := 3000
authProxies := []net.IP{  // can be nil
    net.ParseIP("123.1.1.1"),
    net.ParseIP("123.1.1.2"),
}
server := ipac.NewFakeServer(port, authProxies)


// run in background, returns immediately
server.StartServing() 
client := ipac.NewClient("username", "password", "http://localhost:3000")
proxies, err := client.GetProxies()
server.Stop()

// run in foreground forever
// only stops on Ctrl+C signal
// page can be in browser at http://localhost:3000
server.ServeForever()
```

<!-- CONTRIBUTING -->
## Contributing

Contributions are what make the open source community such an amazing place to be learn, inspire, and create. Any contributions you make are **greatly appreciated**.

1. Fork the Project
2. Create your Feature Branch (`git checkout -b feature/AmazingFeature`)
3. Commit your Changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the Branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request



<!-- LICENSE -->
## License

Distributed under the MIT License. See `LICENSE` for more information.


<!-- CONTACT -->
## Contact

Duong Thanh Hung - [dthung1602@gmail.com](mailto:dthung1602@gmail.com)

Project Link: [https://github.com/dthung1602/instant-proxy-api-client](https://github.com/dthung1602/instant-proxy-api-client)


<!-- ACKNOWLEDGEMENTS -->
## Acknowledgements
* [Best README template](https://github.com/othneildrew/Best-README-Template)
* [Img Shields](https://shields.io)



<!-- MARKDOWN LINKS & IMAGES -->
<!-- https://www.markdownguide.org/basic-syntax/#reference-style-links -->
[contributors-shield]: https://img.shields.io/github/contributors/dthung1602/instant-proxy-api-client.svg?style=flat-square
[contributors-url]: https://github.com/dthung1602/instant-proxy-api-client/graphs/contributors
[forks-shield]: https://img.shields.io/github/forks/dthung1602/instant-proxy-api-client.svg?style=flat-square
[forks-url]: https://github.com/dthung1602/instant-proxy-api-client/network/members
[stars-shield]: https://img.shields.io/github/stars/dthung1602/instant-proxy-api-client.svg?style=flat-square
[stars-url]: https://github.com/dthung1602/instant-proxy-api-client/stargazers
[issues-shield]: https://img.shields.io/github/issues/dthung1602/instant-proxy-api-client.svg?style=flat-square
[issues-url]: https://github.com/dthung1602/instant-proxy-api-client/issues
[license-shield]: https://img.shields.io/github/license/dthung1602/instant-proxy-api-client.svg?style=flat-square
[license-url]: https://github.com/dthung1602/instant-proxy-api-client/blob/master/LICENSE

