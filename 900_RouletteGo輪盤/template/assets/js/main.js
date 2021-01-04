const LEFT = "left";
const RIGHT = "right";

const EVENT_MESSAGE = "message"
const EVENT_OTHER = "other"
const EVENT_ROULETTE = "roulette"

const userPhotos = [
    "https://www.flaticon.com/svg/static/icons/svg/3408/3408584.svg",
    "https://www.flaticon.com/svg/static/icons/svg/3408/3408537.svg",
    "https://www.flaticon.com/svg/static/icons/svg/3408/3408540.svg",
    "https://www.flaticon.com/svg/static/icons/svg/3408/3408545.svg",
    "https://www.flaticon.com/svg/static/icons/svg/3408/3408551.svg",
    "https://www.flaticon.com/svg/static/icons/svg/3408/3408556.svg",
    "https://www.flaticon.com/svg/static/icons/svg/3408/3408564.svg",
    "https://www.flaticon.com/svg/static/icons/svg/3408/3408571.svg",
    "https://www.flaticon.com/svg/static/icons/svg/3408/3408578.svg",
    "https://www.flaticon.com/svg/static/icons/svg/3408/3408720.svg"
]
var PERSON_IMG = userPhotos[getRandomNum(0, userPhotos.length)];
var PERSON_NAME = "Guest" + Math.floor(Math.random() * 1000);

var url = "ws://" + window.location.host + "/ws?id=" + PERSON_NAME;
var ws = new WebSocket(url);
var chatroom = document.getElementsByClassName("msger-chat")
var text = document.getElementById("msg");
var send = document.getElementById("send")

send.onclick = function (e) {
    handleMessageEvent()
}

text.onkeydown = function (e) {
    if (e.keyCode === 13 && text.value !== "") {
        handleMessageEvent()
    }
};

ws.onmessage = function (e) {
    var m = JSON.parse(e.data)
    var msg = ""
    switch (m.event) {
        case EVENT_MESSAGE:
            if (m.name == PERSON_NAME) {
                msg = getMessage(m.name, m.photo, RIGHT, m.content);
            } else {
                msg = getMessage(m.name, m.photo, LEFT, m.content);
            }
            break;
        case EVENT_OTHER:
            if (m.name != PERSON_NAME) {
                msg = getEventMessage(m.name + " " + m.content)
            } else {
                msg = getEventMessage("您已" + m.content)
            }
            break;
        case EVENT_ROULETTE:
            // 這邊是收到後端開幾數字的開獎行為
            msg = getEventMessage("本局輪盤開出號碼為： "+m.content)

            let randomNumber = parseInt(m.content)
            resetFunc()
            spinFunc(randomNumber)
            break;
    }
    insertMsg(msg, chatroom[0]);
};

ws.onclose = function (e) {
    console.log(e)
}

function handleMessageEvent() {
    ws.send(JSON.stringify({
        "event": "message",
        "photo": PERSON_IMG,
        "name": PERSON_NAME,
        "content": text.value,
    }));
    text.value = "";
}

function getEventMessage(msg) {
    var msg = `<div class="msg-left">${msg}</div>`
    return msg
}

function getMessage(name, img, side, text) {
    const d = new Date()
    //   Simple solution for small apps
    var msg = `
    <div class="msg ${side}-msg">
      <div class="msg-img" style="background-image: url(${img})"></div>
      <div class="msg-bubble">
        <div class="msg-info">
          <div class="msg-info-name">${name}</div>
          <div class="msg-info-time">${d.getFullYear()}/${d.getMonth()}/${d.getDay()} ${d.getHours()}:${d.getMinutes()}</div>
        </div>
        <div class="msg-text">${text}</div>
      </div>
    </div>
  `
    return msg;
}

function insertMsg(msg, domObj) {
    domObj.insertAdjacentHTML("beforeend", msg);
    domObj.scrollTop += 500;
}

function getRandomNum(min, max) {
    return Math.floor(Math.random() * (max - min + 1)) + min;
}

function resetFunc(){
    // remove the spinto data attr so the ball 'resets'
    $inner.attr('data-spinto','').removeClass('rest');
    $('#reset').hide();
    $spin.show();
    $data.removeClass('reveal');
}

function spinFunc(randomNumber){

    // 產生隨機數給 nth-child 選擇器 get a random number between 0 and 36 and apply it to the nth-child selector
    //var  randomNumber = Math.floor(Math.random() * 36),
    color = null;
    $inner.attr('data-spinto', randomNumber).find('li:nth-child('+ randomNumber +') input').prop('checked','checked');
    // 隱藏旋轉按鈕可防止重複單擊 prevent repeated clicks on the spin button by hiding it
    $('#spin').hide();
    // 禁用重置按鈕，直到球停止旋轉 disable the reset button until the ball has stopped spinning
    $reset.addClass('disabled').prop('disabled','disabled').show();

    $('.placeholder').remove();

    // 停球時刪除禁用的屬性 remove the disabled attribute when the ball has stopped
    setTimeout(function() {
        $reset.removeClass('disabled').prop('disabled','');

        if($.inArray(randomNumber, red) !== -1){ color = 'red'} else { color = 'black'};
        if(randomNumber == 0){color = 'green'};

        $('.result-number').text(randomNumber);
        $('.result-color').text(color);
        $('.result').css({'background-color': ''+color+''});
        $data.addClass('reveal');
        $inner.addClass('rest');

        $thisResult = '<li class="previous-result color-'+ color +'"><span class="previous-number">'+ randomNumber +'</span><span class="previous-color">'+ color +'</span></li>';
        $('.previous-list').prepend($thisResult);

        // 最多顯示四筆資料
        if ($('.previous-list li').length >= 5){
            console.log("test" +$('.previous-list li').length)
            $('.previous-list li').last().remove();
        }

    }, 9000);
}