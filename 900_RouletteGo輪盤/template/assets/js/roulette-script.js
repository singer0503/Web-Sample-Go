
 var $inner = $('.inner'),
     $spin = $('#spin'),
     $reset = $('#reset'),
     $data = $('.data'),
     $mask = $('.mask'),
     maskDefault = 'Place Your Bets',
     timer = 9000;

var red = [32,19,21,25,34,27,36,30,23,5,16,1,14,9,18,7,12,3];

$reset.hide();

$mask.text(maskDefault);

// $spin.on('click', spinFunc) // Note: 若加上括號在載入時就會執行, 不加上撥號又點像是 IIFE

// $reset.on('click',resetFunt);

// so you can swipe it too
var myElement = document.getElementById('plate');
var mc = new Hammer(myElement);
mc.on("swipe", function(ev) {
  if(!$reset.hasClass('disabled')){
    if($spin.is(':visible')){
      $spin.click();  
    } else {
      $reset.click();
    }
  }  
});

 // IIFE (Immediately Invoked Function Expression)， 一用就丟 立刻被呼叫、執行的 function 表達式。
//  function spinFunc(){
//
//     // 產生隨機數 get a random number between 0 and 36 and apply it to the nth-child selector
//     var randomNumber = Math.floor(Math.random() * 36),
//         color = null;
//     $inner.attr('data-spinto', randomNumber).find('li:nth-child('+ randomNumber +') input').prop('checked','checked');
//     // 隱藏旋轉按鈕可防止重複單擊 prevent repeated clicks on the spin button by hiding it
//     $(this).hide();
//     // 禁用重置按鈕，直到球停止旋轉 disable the reset button until the ball has stopped spinning
//     $reset.addClass('disabled').prop('disabled','disabled').show();
//
//     $('.placeholder').remove();
//
//     console.log('1')
//
//     setTimeout(function() {
//         $mask.text('No More Bets');
//     }, timer/2);
//
//     setTimeout(function() {
//         $mask.text(maskDefault);
//     }, timer+500);
//
//     console.log('2')
//
//     // 停球時刪除禁用的屬性 remove the disabled attribute when the ball has stopped
//     setTimeout(function() {
//         $reset.removeClass('disabled').prop('disabled','');
//
//         if($.inArray(randomNumber, red) !== -1){ color = 'red'} else { color = 'black'};
//         if(randomNumber == 0){color = 'green'};
//
//         $('.result-number').text(randomNumber);
//         $('.result-color').text(color);
//         $('.result').css({'background-color': ''+color+''});
//         $data.addClass('reveal');
//         $inner.addClass('rest');
//
//         $thisResult = '<li class="previous-result color-'+ color +'"><span class="previous-number">'+ randomNumber +'</span><span class="previous-color">'+ color +'</span></li>';
//
//         $('.previous-list').prepend($thisResult);
//
//         console.log('3')
//     }, timer);
//     console.log('4')
// }