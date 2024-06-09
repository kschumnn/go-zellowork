# go-zellowork

## Table of Contents
+ [About](#about)
+ [Usage](#usage)

## About <a name = "about"></a>
I've made this library, since there is no golang implementation available, to manage Zello Work users and channels.

### Prerequisites

You need a valid ZelloWork Instance aswell as an active API Token + User Credentials.

### Installing

`go get github.com/kschumnn/go-zellowork`

## Usage <a name = "usage"></a>
Example for Creating a Channel
```go
import (
    zellowork "github.com/kschumnn/go-zellowork/api"
)
ac := zellowork.NewAPIClient("https://example.zellowork.com", "abcdefg12345678234567")
_, err := ac.Authenticate(username, password)
if err != nil {
    panic(err)
}
err = ac.ChannelAdd("Test Channel", false, false)
if err != nil {
    panic(err)
}
```