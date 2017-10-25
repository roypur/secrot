package secrot

import("fmt"
       "strings"
       "crypto/rand"
       "encoding/hex"
       "sync"
       "os"
       "time"
   )
var secretMutex sync.RWMutex

type Stack struct {
    secretSize uint
    secretCount uint
    secrets []string
    interval time.Duration
    secretMutex sync.RWMutex
}

func (s Stack)SecretAge(str string)(age int, valid bool){
    s.secretMutex.RLock()
    for k,v := range s.secrets{
        if strings.Contains(strings.ToLower(str), v){
            age = k
            valid = true
            break
        }
    }
    s.secretMutex.RUnlock()
    return
}

func (s Stack)GetSecret()(ret string){
    s.secretMutex.RLock()
    ret = s.secrets[0]
    s.secretMutex.RUnlock()
    return
}

func NewStack(secretSize uint, secretCount uint, interval time.Duration)Stack{
    stack := Stack{}
    stack.secretSize = secretSize
    stack.interval = interval
    stack.secrets = make([]string, secretCount)
    initSecrets(stack)
    go updateInterval(stack)

    return stack
}

func updateInterval(s Stack){
    for{
        time.Sleep(s.interval)
        updateSecrets(s)
    }
}

func updateSecrets(s Stack){
    buf := make([]byte, s.secretSize)
    _, err := rand.Read(buf)

    if err != nil{
        fmt.Println("Failed to get random numbers. Crashing...")
        os.Exit(1)
    }
    hexData := hex.EncodeToString(buf)

    s.secretMutex.Lock()
    lastElem := len(s.secrets) - 1
    for i:=lastElem; i > 0; i--{
        s.secrets[i] = s.secrets[i-1]
    }
    s.secrets[0] = hexData
    s.secretMutex.Unlock()
}
func initSecrets(s Stack){
    buf := make([]byte, s.secretSize)
    _, err := rand.Read(buf)

    if err != nil{
        fmt.Println("Failed to get random numbers. Crashing...")
        os.Exit(1)
    }
    hexString := hex.EncodeToString(buf)

    for i,_ := range s.secrets{
        s.secrets[i] = hexString
    }
}
