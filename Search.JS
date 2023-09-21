function displayTime() {
    var time = new Date();
    var hour = time.getHours();
    var minute = time.getMinutes();
    var month = time.getMonth() + 1;
    var day = time.getDate();
    var timeText = document.getElementById("timeText");
    var dateText = document.getElementById("dateText");
    var fhour = hour < 10 ? "0" + hour : hour;
    var fminute = minute < 10 ? "0" + minute : minute;
    var fday = day < 10 ? "0" + day : day;
    timeText.innerHTML = fhour + "<span>:</span>" + fminute;
    dateText.innerHTML = month + "<span>月</span>" + fday + "<span>日</sapn>";
}
window.onload = displayTime;
setInterval(displayTime, 1000);

function addSearchHistory(history) {
    const historyBox = document.getElementById("historyBox");
    const historyItems = document.querySelectorAll(".searchHistory");
    for (item of historyItems) {
        if (item.textContent.trim() == history) {
            item.remove();
        }
    }
    const historyItem = document.createElement("div");
    historyItem.className = "searchHistory";
    historyItem.innerHTML = history;
    historyBox.insertBefore(historyItem, historyBox.firstChild);
    if (historyItems.length > 10) {
        historyItems[0].remove();
    }
}

function getLocalHistory() {
    const historyJson = localStorage.getItem("searchHistory");
    console.log(historyJson);
    if (historyJson) {
        const historyList = JSON.parse(historyJson);
        console.log(historyList);
        historyList.forEach((item) => {
            addSearchHistory(item);
        });
    }
}

function saveLocalHistory() {
    const historyList = document.querySelectorAll(".searchHistory");
    console.log(historyList);
    const historyArray = [];
    historyList.forEach((item) => {
        historyArray.push(item.textContent.trim());
    });
    console.log(historyArray);
    const historyJson = JSON.stringify(historyArray);
    console.log(historyJson);
    localStorage.setItem("searchHistory", historyJson);
}
document.addEventListener("DOMContentLoaded", getLocalHistory);




document.addEventListener("DOMContentLoaded", function () {
    const btnSearch = document.getElementById("btnSearch");
    const btnTrans = document.getElementById("btnTrans");
    const btnBack = document.getElementById("btnBack");
    const btnFire = document.getElementById("btnFire");
    const menuSearchEng = document.getElementById("menuSearchEng");
    const baidu = document.getElementById("baidu");
    const google = document.getElementById("google");
    const bing = document.getElementById("bing");
    const clearHistory = document.getElementById("clearHistory");
    var searchEngine = "https://www.google.com/search?q=";
    var searchInput = document.getElementById("search-input");
    var historyList = document.querySelectorAll(".searchHistory");
    function searchText() {
        var searchText = searchInput.value.trim();
        if (searchText != "") {
            let searchUrl = searchEngine + encodeURIComponent(searchText);
            window.open(searchUrl, "_blank");
            addSearchHistory(searchText);
            saveLocalHistory();
        } else {
            alert("请输入搜索内容！");
        }
    }
    btnSearch.addEventListener("click", () => {
        searchText();
    });
    searchInput.addEventListener("keydown", (e) => {
        if (e.key == "Enter") {
            searchText();
        }
    });
    btnTrans.addEventListener("click", () => {
        if (menuSearchEng.style.display == "none") {
            menuSearchEng.style.display = "block";
        } else {
            menuSearchEng.style.display = "none";
        }
    });
    baidu.addEventListener("click", () => {
        searchEngine = "https://www.baidu.com/s?wd=";
        menuSearchEng.style.display = "none";
    });
    google.addEventListener("click", () => {
        searchEngine = "https://www.google.com/search?q=";
        menuSearchEng.style.display = "none";
    });
    bing.addEventListener("click", () => {
        searchEngine = "https://cn.bing.com/search?q=";
        menuSearchEng.style.display = "none";
    });
    historyList.forEach((Item) => {
        Item.addEventListener("click", () => {
            searchText();
        });
    });
    clearHistory.addEventListener("click", () => {
        const historyList = document.querySelectorAll(".searchHistory");
        historyList.forEach((item) => {
            item.remove();
        });
    });
    btnBack.addEventListener("click", () => {
        searchInput.value = historyList[historyList.length - 1].textContent.trim();
    });
    btnFire.addEventListener("click", () => {
        const four = ["南大家园", "云家园", "家园工作室", "小家园传声机"];
        var random = Math.floor(Math.random() * 4);
        searchInput.value = four[random];
    });
});
