

const maxi = 21;

let sumnum = 1
	, lockweek = 1024 // 线性锁定（周）
    , coin = 1048576
    , total_num = 0
    , total_coin = 0;
for (let i=maxi; i>0; i--) {
    total_coin += coin * sumnum;
    total_num += sumnum;

    console.log("LV: " + prefixZero(maxi-i+1, 2),
        "    BTC: " + prefixZero(sumnum, 7) + ", " + prefixZero(total_num, 7),
		"    HAC: " + prefixZero(coin, 7)   + ", " + prefixZero(total_coin, 8),
		"    LOCK: " + prefixZero(lockweek, 4) + "w, "
			+ prefixZero(((lockweek>0 ? lockweek/52 : 0)+"").substr(0, 5), 5) + "y, "
			+ prefixZero(lockweek>0 ? coin / lockweek : coin, 4)
    );
    sumnum = sumnum * 2;
    coin = coin / 2;
	lockweek = parseInt(lockweek / 2 );
}

let mineryear = 840960;

console.log(`共转移${total_num}枚BTC，增发${total_coin}枚HAC`);












///////////////////////////////////








function prefixZero(num, n) {
    return (Array(n).join(" ") + num).slice(-n);
}

