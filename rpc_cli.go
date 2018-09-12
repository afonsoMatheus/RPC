package main

import (
    "fmt"
    "log"
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



func callLogin(client *rpc.Client, err error)(bool){

	fmt.Println("Digite seu nome de usuário:")
	var user string
	var pass string

	fmt.Scanf("%s", &user)

	fmt.Println("Digite sua senha:")

	fmt.Scanf("%s", &pass)

	////////////

	var reply *rsa.PublicKey

	var cliKey,_ = rsa.GenerateKey(rand.Reader,1024)
	login := User{"",[]byte(""), cliKey.PublicKey}

	err = client.Call("User.ExKey", login, &reply)
	if err != nil {
		log.Fatal("user error:", err)
	}

	crip,_ := rsa.EncryptPKCS1v15(rand.Reader, reply,[]byte(pass))

	login.Name = user
	login.Pass = crip

	////////////

	var reply2 []byte
	err = client.Call("User.Validate", login, &reply2)
	if err != nil {
		log.Fatal("user error:", err)
	}

	var answer,_ = rsa.DecryptPKCS1v15(rand.Reader, cliKey, reply2)


	if string(answer) == "Abandonai toda a esperança vos que entrais"{
		fmt.Println("Você entrou!")
		return true
	}else{
		fmt.Println("Acesso negado!")
		return false
	}

}

func callAdd( item Item ,client *rpc.Client, err error){

	var replyBool bool
	err = client.Call("Item.Add", item, &replyBool)
	if err != nil {
		log.Fatal("item error:", err)
	}

	if replyBool{
		fmt.Println("Produto adicionado!")
	}else{
		fmt.Println("Produto já existe!")
	}
}

func callFind( item Item ,client *rpc.Client, err error){
	var replyBool bool
	err = client.Call("Item.Find", item, &replyBool)
	if err != nil {
		log.Fatal("item error:", err)
	}

	if replyBool{
		fmt.Println("Produto existe!")
	}else{
		fmt.Println("Produto não existe!")
	}

}

func callRemove( item Item ,client *rpc.Client, err error){
	var replyBool bool
	err = client.Call("Item.Remove", item, &replyBool)
	if err != nil {
		log.Fatal("item error:", err)
	}

	if replyBool{
		fmt.Println("Produto removido!")
	}else{
		fmt.Println("Produto não existe!")
	}

}

func callList( item Item ,client *rpc.Client, err error){
	var replyList []Item
	err = client.Call("Item.List", item , &replyList)
	if err != nil {
		log.Fatal("item error:", err)
	}
	fmt.Printf("Item: %v\n", replyList)

}

func callUpdate( item Item ,client *rpc.Client, err error){
	var replyBool bool
	err = client.Call("Item.Update", item ,&replyBool)

	if err != nil {
		log.Fatal("itemUpd error:", err)
	}

	if replyBool{
		fmt.Println("Preço atualizado!")
	}else{
		fmt.Println("Produto não encontrado!")
	}

}


func main() {

    client, err := rpc.Dial("tcp", ":1234")
    if err != nil {
        log.Fatal("dialing:", err)
    }

	for {
		if callLogin(client,err){
			break
		}
	}

	for {
		fmt.Println("|	 Menu	   |")
		fmt.Println("|Adicionar: 1  |")
		fmt.Println("|Listar: 2	   |")
		fmt.Println("|Remover: 3    |")
		fmt.Println("|Procurar: 4   |")
		fmt.Println("|Atualizar: 5  |")
		fmt.Println("|Sair: 0       |")

		var input int
		n, err := fmt.Scanln(&input)

		if n < 1 || err != nil {
			fmt.Println("invalid input")
			return
		}

		switch input {
		case 1:
			fmt.Println("Digite o nome do produto e o seu preço!")

			var name string
			var value int

			fmt.Scanf("%s", &name)
			fmt.Scanln(&value)

			var item = Item{name, value}

			callAdd(item, client, err)
			fmt.Println()

		case 2:
			var item = Item{"Debug", 0}

			callList(item, client, err)
			fmt.Println()

		case 3:

			fmt.Println("Digite o nome do produto que você quer remover!")

			var name string

			fmt.Scanf("%s", &name)

			var item = Item{name, 0}

			callRemove(item, client, err)
			fmt.Println()

		case 4:

			fmt.Println("Digite o nome do produto que você quer procurar!")

			var name string

			fmt.Scanf("%s", &name)

			var item = Item{name, 0}

			callFind(item, client, err)
			fmt.Println()

		case 5:
			fmt.Println("Digite o nome do produto que você quer atualizar!")

			var nam string
			fmt.Scanf("%s", &nam)

			fmt.Println("Digite o novo preço do produto!")

			var value int
			fmt.Scanln(&value)

			var item = Item{nam,value}

			callUpdate(item, client, err)
			fmt.Println()

		case 0:
			os.Exit(2)

		default:
		}
	}









}
