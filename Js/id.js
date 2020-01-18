console.log("ID APPPPP")
let app = new Vue({
    el:"#text",
    data: {
        content: document.getElementById("hide").textContent
    },
    methods: {
        decode: async (event) => {
            let key = document.getElementById("textDecode").textContent
            console.log(key)
        }
    }


})
let form = new Vue({
    el:"#form",
    data: {

    },
    methods: {
        decode: async (event) => {
            let key = document.getElementById("textDecode").value
            let id = document.location.href.split("/")[4]
            console.log(id)
           let data = {key: key, id: id}
            let r = await fetch("/ajax", {
               method: "POST",
               body: JSON.stringify(data)
           })
            let js = await r.json()
             app.content = js.DecodeText



        }
    }


})

