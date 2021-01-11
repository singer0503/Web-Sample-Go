(function($) {
	
	// table
	(function() {
		"use strict"
		
		function getButtonCells(btn) {
			var cells = btn.data('cells');
			if (!cells || !cells.length) {
				cells = [];
				switch (btn.data('type')) {
					case 'sector':
						var nums = sectors[btn.data('sector')];
						for (var i = 0, len = nums.length; i < len; i++) {
							cells.push(table_nums[nums[i]]);
						}
						return cells;
						break;
					case 'num':
					default:
						var nums = String(btn.data('num')).split(',');
						for (var i = 0, len = nums.length; i < len; i++) {
							cells.push(table_nums[nums[i]]);
						}
						btn.data('cells', cells)
						return btn.data('cells');
						break;
				}
			}

			return cells;
		};

		// props
		var active = true,
			selectors = {
				roulette : '.roulette',
				num : '.num',
				sector : '.sector',
				table_btns : '.controlls .btn'
			},
			classes = {
				red : 'red',
				black : 'black',
				green : 'green',
				hover : 'hover'
			},
			numbers = {
				red : [],
				black : [],
				green : []
			},
			sectors = {
				'1' : [], // 1st row
				'2' : [], // 2nd row
				'3' : [], // 3rd row
				'4' : [], // 1st 12
				'5' : [], // 2nd 12
				'6' : [], // 3rd 12
				'7' : [], // 1 to 18
				'8' : [], // EVEN
				'9' : [], // RED
				'10' : [], // BLACK
				'11' : [], // ODD
				'12' : [], // 19 to 36
			},
			table_nums = {},
			table_sectors = {};

		// init
		$(selectors.num).each(function() {
			var $this = $(this),
				color,
				num = Number($this.text());
			// add to instances array
			table_nums[num] = $this;
			// add to colors array
			for (var color in numbers) {
				if ($this.hasClass(classes[color])) {
					numbers[color].push(num);
					$this.data('color', color);
				}
			}
		})

		$(selectors.sector).each(function() { 
			var $this = $(this),
				color;
			if ($this.hasClass(classes.red)) {
				color = 'red';
			} else if ($this.hasClass(classes.black)) {
				color = 'black';
			} else {
				color = 'sector';
			}
			$this.data('color', color);
			table_sectors[$this.data('sector')] = $this;
		});

		// sort numbers
		for (var color in numbers) {
			numbers[color].sort(function(a, b) { return a - b; });
		}

		// populate sectors
		for (var i = 1; i <= 36; i++) {
			// 1st row, 2nd row, 3rd row
			switch (i%3) {
				case 0:
					sectors['1'].push(i);
					break;
				case 1:
					sectors['3'].push(i);
					break;
				case 2:
					sectors['2'].push(i);
					break;
			}

			// 1st 12, 2nd 12, 3rd 12
			if (i <= 12) {
				sectors['4'].push(i);
			} else if (i <= 24) {
				sectors['5'].push(i);
			} else {
				sectors['6'].push(i);
			}

			// 1 to 18, 19 to 36
			if (i <= 18) {
				sectors['7'].push(i);
			} else {
				sectors['12'].push(i);
			}

			// ODD, EVEN
			if (i%2) {
				sectors['11'].push(i);
			} else {
				sectors['8'].push(i);
			}

			if (numbers.red.indexOf(i) != -1) {
				sectors['9'].push(i);
			} else if (numbers.black.indexOf(i) != -1) {
				sectors['10'].push(i);
			}
		}

		// buttons
		var table_btns = $(selectors.table_btns).hover(
			function() {
				hovering=1;
				if (active) {
					var $this = $(this),
						cells = getButtonCells($this);
					for (var i = 0, len = cells.length; i < len; i++) {
						cells[i].addClass(classes.hover);
					}
					var sector = $this.data('sector');
					if (sector) {
						table_sectors[sector].addClass(classes.hover);
					}
				}
			},
			function() {
				hovering=0;
				var $this = $(this),
					cells = getButtonCells($this);
				for (var i = 0, len = cells.length; i < len; i++) {
					cells[i].removeClass(classes.hover);
				}
				var sector = $this.data('sector');
				if (sector) {
					table_sectors[sector].removeClass(classes.hover);
				}
			}
		).mousedown(function(e) {
			var numbers=[];
			if(typeof $(this).data('sector') != 'undefined'){
				console.log("SECTOR "+$(this).data('sector'));
				
				if(e.button==2)ChangeBet(36+$(this).data('sector'),-1);
				else ChangeBet(36+$(this).data('sector'),+1);
			}
			else{
				numbers=$(this).data('num');
				
				if(typeof numbers.length ==='undefined')numbers=[numbers];
				else numbers=numbers.split(',');
				
				if(e.button==2)for(var i=0;i<numbers.length;i++)ChangeBet(numbers[i],-1);
				else for(var i=0;i<numbers.length;i++)ChangeBet(numbers[i],+1);
			}
		});
	})();
	
document.oncontextmenu = function() {if(hovering)return false;}; // 如果滑鼠移動到下注盤上，就把右鍵功能取消

})(jQuery);


var squares=new Array(48);
var divs=document.getElementsByTagName("div");
for(var i=0;i<divs.length;i++){
	var attr=divs[i].getAttribute("data-num");
	if(attr==null){
		attr=divs[i].getAttribute("data-sector");
		if(attr==null)continue;
		var index=36+parseInt(attr);
		console.log(divs[i].getBoundingClientRect()) // TODO: 還是有點偏移量
		var rekt=divs[i].getBoundingClientRect();
		squares[index]=new Array(2);
		squares[index][1]=rekt.top+10;
		squares[index][0]=rekt.left+16;
	}else{
		if(attr.indexOf(',')!=-1)continue;
		var rekt=divs[i].getBoundingClientRect();
		squares[attr]=new Array(2);
		squares[attr][1]=rekt.top+10;
		squares[attr][0]=rekt.left+16;
	}
}

// 繪製 betdiv 下注記錄區塊 動作為 ->清空整個 div 再把已下住的 bets 統計後繪製計算
function UpdateBets(){
	var betdiv=document.getElementById("bets");
	betdiv.innerHTML='';
	for(var i=37;i<bets.length;i++)if(bets[i]>0)betdiv.innerHTML+=sectors[i-37]+": "+(bets[i]*CurrentTier).toFixed(2)+"<br>";
	for(var i=0;i<37;i++)if(bets[i]>0)betdiv.innerHTML+="Number "+i+": "+(bets[i]*CurrentTier).toFixed(2)+"<br>";
}

// 重置所有參數回初始狀態
function Reset(){
	for(var i=0;i<bets.length;i++){
		bets[i]=0;
		if(chips[i]!=null)for(var j=0;chips[i].length>0;j++)document.body.removeChild(chips[i].pop());
	}
	balance=10000; //TODO 要去拿到目前帳戶的餘額
	
	UpdateBets();
	UpdateBalance();
}

//計算所有總下注額
function TotalBets(){
	var r=0;
	for(var i=0;i<bets.length;i++)r+=bets[i];
	return r;
}

// 產生隨機數使用的方法
function rInt(min,max){
	return Math.floor(Math.random() * (max - min + 1)) + min;
}

// 放下注圖片使用的 img array
var chips=new Array(48);

// 將下注圖片繪製於 div 之上
function ChangeBet(id,amount){
	if(amount>0&&TotalBets()==50){
		//maybe some beep
		return;
	}
	
	if(amount>0){
		var img = document.createElement('img');
		img.src="assets/images/bet-chip-100.png";	// 籌碼圖片
		img.style.zIndex="0";
		img.style.position="absolute";
		
		var rX=rInt(-16,16);
		var rY=rInt(-16,16);
		
		img.style.left=(squares[id][0]+rX)+"px";
		img.style.top=(squares[id][1]+rY)+"px";
		
		img.style.width="20px";
		img.style.pointerEvents="none";
		
		document.body.appendChild(img);
		
		if(chips[id]==null)chips[id]=new Array(0);
		chips[id].push(img);
	}if(amount<0&&chips[id]!=null&&chips[id].length>0)document.body.removeChild(chips[id].pop());
	
	bets[id]+=amount;
	if(bets[id]<0)bets[id]=0;
	UpdateBets();
	UpdateBalance();
}

// 更新目前現金資產畫面 餘額及下注總額
function UpdateBalance(){
	var e=document.getElementById("balance");
	e.innerHTML="目前餘額:「 " + balance + " 」 ";
	var tb=TotalBets();
	if(tb>0)e.innerHTML+="本局下注總額 (" + (tb*CurrentTier) + ")";
}

// 開獎
function Place(){
	var bet=0;
	for(var i=0;i<bets.length;i++)if(bets[i]!=0)bet+=bets[i];
	bet*=CurrentTier;
	
	if(bet>balance){
		var betdiv=document.getElementById("result");
		betdiv.innerHTML="餘額不足!!";
		return;
	}
	
	var result=Math.floor(Math.random()*37);
	console.log('本次開獎號碼 result === '+result)
	var win=0;
	if(bets[result]!=0)win+=bets[result]*36; //如果壓中那個號碼是 36 倍的賠率
	for(var i=37;i<bets.length;i++){ //從 bets 37 以上開始算的意思是，0~37 是獨立的數字以在上面的邏輯就以算好，這邊迴圈是計算組合型投注的賠率！
		if(bets[i]!=0){
			console.log('i=== '+i)
			console.log('sectormultipliers[i-37][result] === '+sectormultipliers[i-37][result])
			win+=bets[i]*sectormultipliers[i-37][result]; //計算陪率
		}
	}
	
	win*=CurrentTier; // 計算籌碼價值，目前都是 100 元
	win-=bet; // 減掉投注額，就是贏回來的錢
	
	console.log("下注(bet): "+bet+" 返現(win): "+win);
	
	var betdiv=document.getElementById("result");
	if(bet!=0){ // 不等於 0 才代表有下注，有下注才需要訊息！
		if(win>=bet)betdiv.innerHTML="本局輪盤開出號碼: "+result+" you won " + win + " !!";
		else betdiv.innerHTML="本局輪盤開出號碼: "+result+" you lost " + win + " !!";
	}else{
		betdiv.innerHTML="本局輪盤開出號碼: "+result+" !!";
	}

	
	balance+=win; //持有的錢加上贏來的錢
	UpdateBalance();
}

var balance=10000; // 初始金額
var CurrentTier=100; // 目前設定，每一個籌碼的價值為 100 元
// 應該會有不同金額的籌碼設定，但目前沒用到 (ETH)
var tiers=[
	0.0001,
	0.0002,
	0.001,
	0.002,
	0.01,
	0.02
];

var sectors=[
	"3rd column",
	"2nd column",
	"1st column",
	"1st 12",
	"2nd 12",
	"3rd 12",
	"1 to 18",
	"Even",
	"Red",
	"Black",
	"Odd",
	"19 to 36"
];

var hovering=0;

// 計算下注在哪個區域使用的變數
var bets=[0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0];

// 每個位置的賠率, 二維陣列表
var sectormultipliers=[
	[0,0,0,3,0,0,3,0,0,3,0,0,3,0,0,3,0,0,3,0,0,3,0,0,3,0,0,3,0,0,3,0,0,3,0,0,3],//3rd column
	[0,0,3,0,0,3,0,0,3,0,0,3,0,0,3,0,0,3,0,0,3,0,0,3,0,0,3,0,0,3,0,0,3,0,0,3,0],//2nd column
	[0,3,0,0,3,0,0,3,0,0,3,0,0,3,0,0,3,0,0,3,0,0,3,0,0,3,0,0,3,0,0,3,0,0,3,0,0],//1st column
	[0,3,3,3,3,3,3,3,3,3,3,3,3,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0],//1st 12
	[0,0,0,0,0,0,0,0,0,0,0,0,0,3,3,3,3,3,3,3,3,3,3,3,3,0,0,0,0,0,0,0,0,0,0,0,0],//2nd 12
	[0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,3,3,3,3,3,3,3,3,3,3,3,3],//3rd 12
	[0,2,2,2,2,2,2,2,2,2,2,2,2,2,2,2,2,2,2,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0],//1 to 18
	[0,0,2,0,2,0,2,0,2,0,2,0,2,0,2,0,2,0,2,0,2,0,2,0,2,0,2,0,2,0,2,0,2,0,2,0,2],//even
	[0,2,0,2,0,2,0,2,0,2,0,0,2,0,2,0,2,0,2,2,0,2,0,2,0,2,0,2,0,0,2,0,2,0,2,0,2],//Red
	[0,0,2,0,2,0,2,0,2,0,2,2,0,2,0,2,0,2,0,0,2,0,2,0,2,0,2,0,2,2,0,2,0,2,0,2,0],//Black
	[0,2,0,2,0,2,0,2,0,2,0,2,0,2,0,2,0,2,0,2,0,2,0,2,0,2,0,2,0,2,0,2,0,2,0,2,0],//odd
	[0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,2,2,2,2,2,2,2,2,2,2,2,2,2,2,2,2,2,2] //19 to 36
];