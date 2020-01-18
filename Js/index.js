console.log("Hi")

let app = new Vue({
    el:"#app",
    data: {
        textArea: "Введите свой секретный текст для шифрования"
    },
    methods: {
            sendData:  async(event)=> {
                alert("Send////")
                console.log("SendData")
            let data = {content: app.textArea}
            let req =  await fetch("/data", {
                method: "POST",
                body: JSON.stringify(data)

            })
                let js = await req.json()
                alert("Ваш ключ шифрования: "+js.Key +" сохраните его")
                document.location.href = js.Link


        }
    }
})

