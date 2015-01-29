package main

import "regexp"

func init() {
}

/*
[REL 2015-01-28 09:51:02] 玩家[16266],通过[装备升级],[消耗][银两],数量为[185220],获得[0],数量为[1], 玩家银两数为[12528041] tag:logsvr-1]
[REL 2015-01-28 09:51:02] 玩家[19994],通过[装备升级],[消耗][银两],数量为[177900],获得[0],数量为[1], 玩家银两数为[85524] tag:logsvr-1]
[REL 2015-01-28 09:51:03] 玩家[15365],通过[装备升级],[消耗][银两],数量为[38700],获得[0],数量为[1], 玩家银两数为[8336509] tag:logsvr-1]
[REL 2015-01-28 09:51:04] 玩家[25221],通过[装备升级],[消耗][银两],数量为[1080],获得[0],数量为[1], 玩家银两数为[88790] tag:logsvr-1]
*/
//time0
//player1
//method
//silver
//wupin
//wupin-count
//silver-total
func init() {
	WordDecoder([]string{"装备升级],[消耗][银两],数量为["}, log_equipment_upgrade)
}

var leu = regexp.MustCompile(`\[REL (\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2})\] 玩家\[(\d+)\],通过\[([^\]]+)\],\[[^\]]+\]\[[^\]]+\],数量为\[([^\]]+)\],获得\[([^\]]+)\],数量为\[([^\]]+)\], 玩家银两数为\[([^\]]+)\]`)

func log_equipment_upgrade(line string) doc {
	fields := leu.FindStringSubmatch(line)
	if len(fields) < 8 {
		return nil
	}
	return doc{
		"time":      firstgame_time(fields[1]),
		"player":    iconvert2(fields[2]),
		"type":      "player-" + fields[3],
		"consumes":  iconvert2(fields[4]),
		"wupin":     iconvert2(fields[5]),
		"wupin_cnt": iconvert2(fields[6]),
		"silvers":   iconvert2(fields[7]),
	}
}

/*
[REL 2015-01-28 09:51:02] 玩家[16266][强化]装备[31003-3]从等级[65]到等级[68] tag:logsvr-1 hostname:localhost]
[REL 2015-01-28 09:51:02] 玩家[19994][强化]装备[23009-99]从等级[78]到等级[81] tag:logsvr-1 hostname:localhost]
[REL 2015-01-28 09:51:03] 玩家[15365][强化]装备[23002-79]从等级[47]到等级[50] tag:logsvr-1 hostname:localhost]
[REL 2015-01-28 09:51:04] 玩家[25221][强化]装备[11008-0]从等级[7]到等级[9] tag:logsvr-1 hostname:localhost]
*/
var lee = regexp.MustCompile(`\[REL (\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2})\] 玩家\[(\d+)\]\[强化\]装备\[([^\]]+)\]从等级\[(\d+)\]到等级\[(\d+)\]`)

func init() {
	WordDecoder([]string{"强化]装备[", "从等级[", "到等级["}, log_equipment_enforcement)
}

//time0
//player1
//equipment
//level_from
//level_to
func log_equipment_enforcement(line string) doc {
	fields := lee.FindStringSubmatch(line)
	if len(fields) < 6 {
		return nil
	}
	return doc{
		"time":       firstgame_time(fields[1]),
		"player":     iconvert2(fields[2]),
		"type":       "player-equip",
		"equipment":  fields[3],
		"level_from": iconvert2(fields[4]),
		"level_to":   iconvert2(fields[5]),
	}
}

/*
[REL 2015-01-28 09:51:02] 玩家[14460]通过[卡牌培养],[消耗]道具[200001],数量为[5] tag:logsvr-1 hostname:localhost]
[REL 2015-01-28 09:51:03] 玩家[14058]通过[卡牌培养],[消耗]道具[200001],数量为[5] tag:logsvr-1 hostname:localhost]
[REL 2015-01-28 09:51:04] 玩家[14460]通过[卡牌培养],[消耗]道具[200001],数量为[5] tag:logsvr-1 hostname:localhost]
*/
func init() {
	WordDecoder([]string{"通过[卡牌培养],[消耗]道具["}, log_card_training)
}

var lct = regexp.MustCompile(`\[REL (\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2})\] 玩家\[(\d+)\]通过\[卡牌培养\],\[消耗\]道具\[(\d+)\],数量为\[(\d+)\]`)

//time0
//player1
//equip2
//cnt
func log_card_training(line string) doc {
	fields := lct.FindStringSubmatch(line)
	if len(fields) < 5 {
		return nil
	}
	return doc{
		"time":      firstgame_time(fields[1]),
		"player":    iconvert2(fields[2]),
		"type":      "player-card",
		"equipment": fields[3],
		"cnt":       iconvert2(fields[4]),
	}
}

/*
[REL 2015-01-28 09:51:04] 玩家[14460],通过[卡牌培养],[消耗][元宝],数量为[1],获得[0],数量为[0], 玩家元宝数为[752] tag:logsvr-1]
[REL 2015-01-28 09:51:03] 玩家[14058],通过[卡牌培养],[消耗][元宝],数量为[1],获得[0],数量为[0], 玩家元宝数为[572] tag:logsvr-1 hostname:localhost]
[REL 2015-01-28 09:51:02] 玩家[14460],通过[卡牌培养],[消耗][元宝],数量为[1],获得[0],数量为[0], 玩家元宝数为[753] tag:logsvr-1 hostname:localhost]
*/
func init() {
	WordDecoder([]string{"],通过[卡牌培养],[消耗][元宝],数量为[", "], 玩家元宝数为["}, log_gold_consume)
}

var lgc = regexp.MustCompile(`\[REL (\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2})\] 玩家\[(\d+)\],通过\[卡牌培养\],\[消耗\]\[元宝\],数量为\[(\d+)\],获得\[(\d+)\],数量为\[(\d+)\], 玩家元宝数为\[(\d+)\]`)

//time0, player1
//gold_consume
//equip
//equit-cnt
//gold-left
func log_gold_consume(line string) doc {
	fields := lgc.FindStringSubmatch(line)
	if len(fields) < 7 {
		return nil
	}
	return doc{
		"time":         firstgame_time(fields[1]),
		"player":       iconvert2(fields[2]),
		"gold_consume": iconvert2(fields[3]),
		"equip":        fields[4],
		"equip_cnt":    iconvert2(fields[5]),
		"gold_left":    fields[6],
	}
}

/*
[REL 2015-01-28 09:51:02] 玩家[14460]的[4816]通过[精心培养一次]消耗培养丹[5]培养结果血[-3]攻[0]防[0]内[7]]
[REL 2015-01-28 09:51:03] 玩家[14058]的[4532]通过[精心培养一次]消耗培养丹[5]培养结果血[-1]攻[0]防[0]内[4] tag:logsvr-1 hostname:localhost]
*/
func init() {
	WordDecoder([]string{"培养结果血[", "消耗培养丹["}, log_training)
}

var lt = regexp.MustCompile(`\[REL (\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2})\] 玩家\[(\d+)\]的\[(\d+)\]通过\[([^\]]+)\]消耗培养丹\[(\d+)\]培养结果血\[(-?\d+)\]攻\[(-?\d+)\]防\[(-?\d+)\]内\[(-?\d+)\]`)

func log_training(line string) doc {
	fields := lt.FindStringSubmatch(line)
	if len(fields) < 10 {
		return nil
	}
	return doc{
		"time":    firstgame_time(fields[1]),
		"player":  iconvert2(fields[2]),
		"type":    "player-training",
		"card":    iconvert2(fields[3]),
		"method":  fields[4],
		"dan":     iconvert2(fields[5]),
		"hp":      iconvert2(fields[6]),
		"force":   iconvert2(fields[7]),
		"defence": iconvert2(fields[8]),
		"qi":      iconvert2(fields[9]),
	}
}

/*
[REL 2015-01-28 09:51:03] 玩家[24206]将伙伴[1506]上阵，换下[-1]，上阵数量[2-2] tag:logsvr-1]
*/
func init() {
	WordDecoder([]string{"]将伙伴[", "]上阵，换下[", "]，上阵数量["}, log_training)
}

var lp = regexp.MustCompile(`\[REL (\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2})\] 玩家\[(\d+)\]将伙伴\[(\d+)\]上阵，换下\[(-?\d+)\]，上阵数量\[([^\]]+)\]`)

func log_partner(line string) doc {
	fields := lp.FindStringSubmatch(line)
	if len(fields) < 6 {
		return nil
	}
	return doc{
		"time":      firstgame_time(fields[1]),
		"player":    iconvert2(fields[2]),
		"type":      "player-partner",
		"partner":   iconvert2(fields[3]),
		"exit_cnt":  iconvert2(fields[4]),
		"enter_cnt": iconvert2(fields[5]),
	}
}

/*
[REL 2015-01-28 09:51:03] 玩家[12993],获得[9]星,[8]人阵通过闯关[48]关,难度[1],当前战斗力为[171977] tag:logsvr-1 hostname:localhost]
[REL 2015-01-28 09:51:03] 玩家[22372],获得[9]星,[6]人阵通过闯关[11]关,难度[1],当前战斗力为[33487] tag:logsvr-1 hostname:localhost]
[REL 2015-01-28 09:51:03] 玩家[12108],获得[9]星,[8]人阵通过闯关[34]关,难度[1],当前战斗力为[174885] tag:logsvr-1 hostname:localhost]
[REL 2015-01-28 09:51:03] 玩家[23985],获得[6]星,[8]人阵通过闯关[33]关,难度[1],当前战斗力为[112479] tag:logsvr-1 hostname:localhost]
*/
func init() {
	WordDecoder([]string{"人阵通过闯关[", "]关,难度[", "],当前战斗力为["}, log_chuangguan)
}

var lc = regexp.MustCompile(`\[REL (\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2})\] 玩家\[(\d+)\],获得\[(\d+)\]星,\[(\d+)\]人阵通过闯关\[(\d+)\]关,难度\[(\d+)\],当前战斗力为\[(\d+)\]`)

func log_chuangguan(line string) doc {
	fields := lp.FindStringSubmatch(line)
	if len(fields) < 8 {
		return nil
	}
	return doc{
		"time":        firstgame_time(fields[1]),
		"player":      iconvert2(fields[2]),
		"type":        "player-chuanguan",
		"star":        iconvert2(fields[3]),
		"partner_cnt": iconvert2(fields[4]),
		"checkpoint":  iconvert2(fields[5]),
		"difficulty":  iconvert2(fields[6]),
		"force":       iconvert2(fields[7]),
	}
}

func log_ignore(string) doc {
	return doc{}
}
