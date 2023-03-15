// global array
let blogs = []

// function get data from input text html
function getData(e) {
    e.preventDefault()

    // delaclaration variable dom selection
    let title = document.getElementById('input-blog-title').value
    let content = document.getElementById('input-blog-content').value
    let image = document.getElementById('input-blog-image').files

    // Convert spesific image to blob object
    image = URL.createObjectURL(image[0])

    let dataBlog = {
        title,
        content,
        image,
        author: "Dandi Saputra",
        postedAt: new Date()
    }
    blogs.push(dataBlog)
    console.log(blogs)
    showData()
}

// declaration function show list preview data blog
function showData(){
    document.getElementById("contents").innerHTML = ""
    for(let i=0; i< blogs.length; i++){
        document.getElementById("contents").innerHTML += `
        <div class="blog-list-item">
            <div class="blog-image">
                <img src="${blogs[i].image}" alt="" />
            </div>
            <div class="blog-content">
            <div class="btn-group">
                <button class="btn-edit">Edit Post</button>
                <button class="btn-post">Post Blog</button>
            </div>
            <h1>
                <a href="blog-detail.html" target="_blank">${blogs[i].title}</a>
            </h1>
            <div class="detail-blog-content">
                ${createTime(blogs[i].postedAt)} | ${blogs[i].author}
            </div>
            <p>${blogs[i].content}</p>

            <div style="float:right; margin: 10px">
                <p style="font-size: 15px; color:grey">${getDuration(blogs[i].postedAt)}</p>
            </div>

            </div>
        </div>
        `
    }
}

// function manipulation date tine
function createTime(time){
    // declaration variable
    let years = time.getFullYear()
    let monthIndex = time.getMonth()
    let date = time.getDate()
    let hour = time.getHours()
    let minutes = time.getMinutes()
    
    const month = ["January", "February", "March", "April", "May", "June", "July", "August", "September", "October", "November", "December"]

    return `${date} ${month[monthIndex]} ${years} ${hour}:${minutes} WIB`
}


// manipulation duration time
const getDuration = (timePost) => {
    let timePosting = timePost
    let timeNow = new Date()
    let distance = timeNow - timePosting

    // 1 sec = 1000 milisecond
    // 30*24*60*60*1000

    let monthDistance = Math.floor(distance / (30*24*60*60*1000))
    if(monthDistance > 0) {
        return monthDistance + ' month ago'
    } else {
        let weekDistance = Math.floor(distance / (7*24*60*60*1000))
        if(weekDistance > 0) {
            return weekDistance + ' week ago'
        } else {
            let dayDistance = Math.floor(distance / (24*60*60*1000))
            if(dayDistance > 0) {
                return dayDistance + "Day ago"
            } else {
                let hourDistance = Math.floor(distance / (60*60*1000))
                if(hourDistance > 0) {
                    return hourDistance + "Hour ago"
                } else {
                    let minutesDistance = Math.floor(distance / (60*1000))
                    if(minutesDistance > 0) {
                        return minutesDistance + "Minute ago"
                    } else {
                        let secondDistance = Math.floor(distance / 1000)
                        if (secondDistance > 0) {
                            return secondDistance + "second ago"
                        }
                    }
                }
            }
        }
    }
}

// setInterval(showData, 1000)