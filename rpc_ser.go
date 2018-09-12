package main

import (
    "fmt"
    "net"
    "net/rpc"
    "os"
	"crypto/rsa"
	"crypto/rand"
)

type Item struct{
	Name string
	Value int
}

type User struct{
	Name string
	Pass []byte
	Key rsa.PublicKey

}

var registered = make([]User,100)

var deposit = make([]Item, 100)

var servKey rsa.PrivateKey

func registerUsers()  {

	var debug,_= rsa.GenerateKey(rand.Reader,1024)

	user3 := User{"Ruth",[]byte("ruiva123"), debug.PublicKey}

	registered = append(registered, user3)

}

func (t *User) ExKey(args *User,reply *rsa.PublicKey) error{

	registered[0].Key = args.Key
	*reply = servKey.PublicKey

	return nil

}

func (t *User) GetClientKey(args *User,reply *bool) error{

	registered[0].Key = args.Key

	*reply = true

	return nil

}

func (t *User) Validate(args *User, reply *[]byte) error{

	args.Pass,_ = rsa.DecryptPKCS1v15(rand.Reader, &servKey, args.Pass)

	for i:=0;i < len(registered); i++ {
		if registered[i].Name == args.Name && string(registered[i].Pass) == string(args.Pass) {
			crip,_ := rsa.EncryptPKCS1v15(rand.Reader, &registered[0].Key, []byte("Abandonai toda a esperança vos que entrais"))
			*reply = crip

			return nil
		}
	}

	crip,_ := rsa.EncryptPKCS1v15(rand.Reader, &registered[0].Key, []byte("supercalifragilisticoespialidoso"))
	*reply = crip
	return nil

}

func (t *Item) Add(args *Item, reply *bool) error{
	for i:=0; i < len(deposit) ; i++ {
		if args.Name == deposit[i].Name {
			*reply = false
			return nil
		}
	}
	deposit = append(deposit, *args)
	*reply = true
	return nil
}

func (t *Item) Remove(args *Item, reply *bool) error{
	for i:=0; i<len(deposit) ; i++ {
		if (args.Name == deposit[i].Name) {
			deposit = append(deposit[:i],deposit[i+1:]...)
			*reply = true
			return nil
		}
	}
	*reply = false
	return nil
}

func (t *Item) Find(args *Item, reply *bool) error{
	for i := 0; i < len(deposit); i++ {
		if(args.Name == deposit[i].Name){
			*reply = true
			return nil
		}
	}
	*reply = false
	return nil
}

func (t *Item) List(args *Item ,reply *[]Item) error{
	for i := 0; i < len(deposit); i++ {
		if deposit[i].Name != "" {
			*reply = append(*reply, deposit[i])
		}
	}
	return nil
}

func (t *Item) Update(args *Item ,reply *bool) error{
	for i := 0; i < len(deposit); i++ {
		if deposit[i].Name == args.Name {
			deposit[i].Value = args.Value
			*reply = true
			return nil
		}
	}
	*reply = false
	return nil
}

func main() {

    item := new(Item)
	user := new(User)

    rpc.Register(item)
	rpc.Register(user)

	var keyS,_= rsa.GenerateKey(rand.Reader,1024)

	servKey = *keyS

	registerUsers()

	tcpAddr, err := net.ResolveTCPAddr("tcp", ":1234")
    checkError(err)

    listener, err := net.ListenTCP("tcp", tcpAddr)
    checkError(err)

    for {
        conn, err := listener.Accept()
        if err != nil {
            continue
        }
        rpc.ServeConn(conn)
    }

}

func checkError(err error) {
    if err != nil {
        fmt.Println("Fatal error ", err.Error())
        os.Exit(1)
    }
}
