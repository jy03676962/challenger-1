//
//  SecondViewController.swift
//  admin
//
//  Created by tassar on 4/20/16.
//  Copyright © 2016 pulupulu. All rights reserved.
//

import UIKit
import Alamofire
import AlamofireImage
import EasyPeasy
import ObjectMapper
import SwiftyJSON
import SwiftyUserDefaults

let cellSize = 45
let cellBorder = 10

class MatchController: PLViewController {
	@IBOutlet weak var groupIDLabel: UILabel!
	@IBOutlet weak var matchStatusLabel: UILabel!
	@IBOutlet weak var playerCountLabel: UILabel!
	@IBOutlet weak var totalCoinLabel: UILabel!
	@IBOutlet weak var energyLabel: UILabel!
	@IBOutlet weak var matchTimeLabel: UILabel!
	@IBOutlet weak var matchModeImageView: UIImageView!
	@IBOutlet weak var mapContainerView: UIView!
	@IBOutlet weak var playerTableView: UITableView!

	var match: Match?

	var mapView: UIImageView = UIImageView()

	@IBAction func forceEnd() {
		let json = JSON([
			"cmd": "stopMatch",
			"matchID": Defaults[.matchID]
		])
		WsClient.singleton.sendJSON(json)
	}

	override func viewDidLoad() {
		super.viewDidLoad()
		playerTableView.backgroundColor = UIColor.clearColor()
	}

	override func viewWillAppear(animated: Bool) {
		super.viewWillAppear(animated)
		DataManager.singleton.subscribeData([.UpdateMatch, .MatchStop], receiver: self)
		if mapView.image == nil {
			Alamofire.request(.GET, PLConstants.getHttpAddress("api/asset/map.png"))
				.responseImage(completionHandler: { response in
					if let image = response.result.value {
						self.mapView.image = image
						self.mapContainerView.addSubview(self.mapView)
						self.mapView <- [
							Size(image.size),
							Center()
						]
					}
			})
		}
	}

	override func viewDidDisappear(animated: Bool) {
		super.viewDidDisappear(animated)
		DataManager.singleton.unsubscribe(self)
	}

	func renderMatch() {
		if match != nil && match!.id == Defaults[.matchID] {
			groupIDLabel.text = match!.teamID
			matchModeImageView.image = match!.mode == "g" ? UIImage(named: "FunIcon") : UIImage(named: "SurvivalIcon")
			let min = Int(match!.elasped) / 60
			let sec = Int(match!.elasped) % 60
			matchTimeLabel.text = String(format: "%02d:%02d", min, sec)
			matchStatusLabel.text = "实时状态: 进行中"
			playerCountLabel.text = "玩家人数:\(match!.member.count)"
			totalCoinLabel.text = "总金币:\(match!.gold)G"
			energyLabel.text = String(format: "%.1f/%d", match!.energy, match!.maxEnergy)
			playerTableView.reloadData()
		} else {
			matchTimeLabel.text = "00: 00"
			matchStatusLabel.text = "实时状态: 未进行"
			playerCountLabel.text = "玩家人数: 0"
			totalCoinLabel.text = "总金币:0G"
			energyLabel.text = ""
			playerTableView.reloadData()
		}
	}
}

extension MatchController: DataReceiver {
	func onReceivedData(json: [String: AnyObject], type: DataType) {
		if type == .UpdateMatch {
			match = Mapper<Match>().map(json["data"] as! String)
			if match != nil && match?.id == Defaults[.matchID] {
				renderMatch()
			}
		} else if type == .MatchStop {
			let matchID = json["data"] as! Int
			if matchID == Defaults[.matchID] {
				match = nil
				renderMatch()
			}
		}
	}
}

extension MatchController: UITableViewDelegate, UITableViewDataSource {
	func tableView(tableView: UITableView, numberOfRowsInSection section: Int) -> Int {
		if match == nil {
			return 0
		} else {
			return match!.member.count
		}
	}

	func numberOfSectionsInTableView(tableView: UITableView) -> Int {
		return 1
	}

	func tableView(tableView: UITableView, cellForRowAtIndexPath indexPath: NSIndexPath) -> UITableViewCell {
		let cell = tableView.dequeueReusableCellWithIdentifier("PlayerTableViewCell") as! PlayerTableViewCell
		cell.setData(match!.member[indexPath.row])
		return cell
	}
}
