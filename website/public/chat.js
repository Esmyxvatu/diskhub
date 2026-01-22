document.getElementById("content").addEventListener("keypress", (event) => {
    if (event.key === "Enter") {
        if (event.shiftKey) return;

        event.preventDefault();
        document.getElementById("form").classList.add("active");

        console.log("Starting request...");

        body = {
            "model": document.getElementById("model").value,
            "content": document.getElementById("content").value,
            "stream": false
        };

        article = document.createElement("article");
        p = document.createElement("p");
        p.innerText += document.getElementById("content").value;
        article.appendChild(p);
        article.className = "message_user";
        document.getElementById("messages").appendChild(article);

        document.getElementById("content").value = "";

        placeHolder = document.createElement("article");
        placeHolder.innerHTML += `
            <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="lucide lucide-loader-circle-icon lucide-loader-circle"><path d="M21 12a9 9 0 1 1-6.219-8.56"/></svg>
            <p>Assistant is thinking ...</p>
        `;
        placeHolder.className = "placeholder";
        document.getElementById("messages").appendChild(placeHolder);

        fetch("/api/askollama", {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify(body)
        }).then(response => response.json()).then(data => {
            console.log(data);

            document.getElementById("messages").removeChild(placeHolder);

            article = document.createElement("article");
            article.innerHTML += `${ data.content }`;
            article.className = `message_${ data.role }`;
            document.getElementById("messages").appendChild(article);
        });
    }
});

if (document.getElementById("messages").children.length > 0) {
    document.getElementById("form").classList.add("active");
}
