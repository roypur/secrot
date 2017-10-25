package secrot

import("fmt"
       "strings"
       "crypto/rand"
       "encoding/hex"
       "sync"
       "os"
       "time"
   )

// This stack contains all the secrets
type Stack struct {
    secretSize uint
    secretCount uint
    secrets []string
    interval time.Duration
    secretMutex sync.RWMutex
}

// Gets the age of a specific secret.
// The newest one has age 0 and so on.
// If valid is false, the secret doesn't exist in the stack.
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

// Gets the newest secret from the stack.
func (s Stack)GetSecret()(ret string){
    s.secretMutex.RLock()
    ret = s.secrets[0]
    s.secretMutex.RUnlock()
    return
}

// Creates a new stack with secrets. secretSize is the size of the secrets
// and secretCount is the number of secrets to keep.
// interval is how often new secrets should be made.
// When a new secret has been made the oldest one is deleted.
func NewStack(secretSize uint, secretCount uint, interval time.Duration)Stack{
    stack := Stack{}
    stack.secretSize = secretSize
    stack.interval = interval
    stack.secrets = make([]string, secretCount)
    initSecrets(stack)
    go updateInterval(stack)

    return stack
}

// Creates a new secret every s.interval.
func updateInterval(s Stack){
    for{
        time.Sleep(s.interval)
        updateSecrets(s)
    }
}

// Pushes the stack and puts a new secret on top.
// The oldest secret will be deleted.
// If it fails to get random numbers
// the program will exit with status code 1.
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

// Creates a new secret and populates the stack with it.
// If it fails to get random numbers
// the program will exit with status code 1.
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
