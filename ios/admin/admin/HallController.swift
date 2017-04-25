//
//  HallController.swift
//  admin
//
//  Created by tassar on 4/20/16.
//  Copyright © 2016 pulupulu. All rights reserved.
//

import UIKit
import SwiftyJSON
import SWTableViewCell
import ObjectMapper
import PKHUD

class HallController: PLViewController {
	fileprivate static let controllerButtonTagStart = 100

	@IBOutlet weak var teamtableView: UITableView!
	@IBOutlet weak var teamIDLabel: UILabel!
	@IBOutlet var controllerButtons: [DeviceButton]!

	@IBOutlet weak var modeImageView: UIImageView!
	@IBOutlet weak var modeLabel: UILabel!
	@IBOutlet weak var playerNumberLabel: UILabel!
	@IBOutlet weak var readyButton: UIButton!
	@IBOutlet weak var startButton: UIButton!
	@IBOutlet weak var addPlayerButton: UIButton!
	@IBOutlet weak var removePlayerButton: UIButton!
	@IBOutlet var changeModeTGR: UITapGestureRecognizer!
	@IBOutlet weak var callButton: UIButton!
	@IBOutlet weak var delayButton: UIButton!

	var refreshControl: UIRefreshControl!
	var teams: [Team]?
	var topTeam: Team?
	var controllers: [PlayerController]?
	var hasPlayingTeam = false

	override func viewDidLoad() {
		super.viewDidLoad()
		refreshControl = UIRefreshControl()
		refreshControl.addTarget(self, action: #selector(HallController.refreshTeamData), for: UIControlEvents.valueChanged)
		teamtableView.addSubview(refreshControl)
	}
	override func viewWillAppear(_ animated: Bool) {
		super.viewWillAppear(animated)
		DataManager.singleton.subscribeData([.HallData, .ControllerData, .NewMatch, .Error], receiver: self)
	}

	override func viewDidDisappear(_ animated: Bool) {
		super.viewDidDisappear(animated)
		DataManager.singleton.unsubscribe(self)
	}
	func refreshTeamData() {
		DataManager.singleton.queryData(.HallData)
	}
	@IBAction func changeMode() {
		guard let team = topTeam else {
			return
		}
		let mode = team.mode == "g" ? "s" : "g"
		let json = JSON([
			"cmd": "teamChangeMode",
			"teamID": topTeam!.id,
			"mode": mode
		])
		WsClient.singleton.sendJSON(json)
	}
	@IBAction func callTeam(_ sender: UIButton) {
		guard let team = topTeam else {
			return
		}
		let json = JSON([
			"cmd": "teamCall",
			"teamID": team.id,
		])
		WsClient.singleton.sendJSON(json)
	}
	@IBAction func delayTeam(_ sender: UIButton) {
		guard let team = topTeam else {
			return
		}
		guard team.status == .waiting else {
			return
		}
		let json = JSON([
			"cmd": "teamDelay",
			"teamID": team.id,
		])
		WsClient.singleton.sendJSON(json)
	}
	@IBAction func addPlayer(_ sender: UIButton) {
		guard let team = topTeam, team.size < PLConstants.maxTeamSize else {
			return
		}
		let json = JSON([
			"cmd": "teamAddPlayer",
			"teamID": team.id,
		])
		WsClient.singleton.sendJSON(json)
		sender.isEnabled = false
	}
	@IBAction func removePlayer(_ sender: UIButton) {
		guard topTeam != nil && topTeam!.size > 1 else {
			return
		}
		let json = JSON([
			"cmd": "teamRemovePlayer",
			"teamID": topTeam!.id,
		])
		WsClient.singleton.sendJSON(json)
		sender.isEnabled = false
	}
	@IBAction func ready(_ sender: UIButton) {
		guard topTeam != nil else {
			return
		}
		if topTeam!.status == .prepare {
			let json = JSON([
				"cmd": "teamCancelPrepare",
				"teamID": topTeam!.id,
			])
			WsClient.singleton.sendJSON(json)
		} else {
			let json = JSON([
				"cmd": "teamPrepare",
				"teamID": topTeam!.id,
			])
			WsClient.singleton.sendJSON(json)
		}
	}
	@IBAction func start(_ sender: AnyObject) {
		guard topTeam != nil else {
			return
		}
		var selectedControllerIds = [String]()
		for btn in controllerButtons {
			guard let pc = btn.controller, btn.isSelected else {
				continue
			}
			selectedControllerIds.append(pc.id)
		}

		let json = JSON([
			"cmd": "teamStart",
			"teamID": topTeam!.id,
			"mode": topTeam!.mode,
			"ids": selectedControllerIds.joined(separator: ",")
		])
		WsClient.singleton.sendJSON(json)
		HUD.show(.progress)
	}

	@IBAction func toggleControllerButton(_ sender: UIButton) {
		sender.isSelected = !sender.isSelected
		startButton.isEnabled = canStart()
	}

	fileprivate func renderTopWaitingTeam() {
		guard topTeam != nil else {
			callButton.isEnabled = false
			delayButton.isEnabled = false
			changeModeTGR.isEnabled = false
			addPlayerButton.isEnabled = false
			removePlayerButton.isEnabled = false
			startButton.isEnabled = false
			readyButton.isEnabled = false
			return
		}
		readyButton.isEnabled = true
		teamIDLabel.text = topTeam!.id
		playerNumberLabel.text = "\(topTeam!.size)"
		if topTeam!.mode == "g" {
			modeImageView.image = UIImage(named: "FunIcon")
			modeLabel.text = "[赏金模式]"
		} else {
			modeImageView.image = UIImage(named: "SurvivalIcon")
			modeLabel.text = "[生存模式]"
		}
		if topTeam!.status == .waiting {
			readyButton.setBackgroundImage(UIImage(named: "PrepareButton"), for: UIControlState())
			callButton.isEnabled = true
			delayButton.isEnabled = true
			changeModeTGR.isEnabled = true
			addPlayerButton.isEnabled = topTeam!.size < PLConstants.maxTeamSize
			removePlayerButton.isEnabled = topTeam!.size > 1
		} else if topTeam!.status == .prepare {
			readyButton.setBackgroundImage(UIImage(named: "CancelPrepare"), for: UIControlState())
			callButton.isEnabled = false
			delayButton.isEnabled = false
			changeModeTGR.isEnabled = false
			addPlayerButton.isEnabled = false
			removePlayerButton.isEnabled = false
		}
		startButton.isEnabled = canStart()
	}

	fileprivate func canStart() -> Bool {
		var count = 0
		for btn in controllerButtons {
			if btn.isSelected {
				count += 1
			}
		}
		return topTeam != nil && topTeam!.status == .prepare && topTeam!.size == count && !hasPlayingTeam
	}

	fileprivate func getBtn(_ idx: Int) -> UIButton {
		return view.viewWithTag(idx + HallController.controllerButtonTagStart) as! UIButton
	}
}

extension HallController: DataReceiver {
	func onReceivedData(_ json: [String: Any], type: DataType) {
		if type == .HallData {
			topTeam = nil
            teams = Mapper<Team>().mapArray(JSONObject: json["data"])
			if teams != nil {
				var topTeamSet = false
				hasPlayingTeam = false
				for team in teams! {
					if (team.status == .waiting || team.status == .prepare) && !topTeamSet {
						topTeam = team
						topTeamSet = true
					}
					if team.status == .playing {
						hasPlayingTeam = true
					}
				}
			}
			renderTopWaitingTeam()
			teamtableView.reloadData()
			refreshControl.endRefreshing()
		} else if type == .ControllerData {
            guard let controllers = Mapper<PlayerController>().mapArray(JSONObject: json["data"]) else {
				return
			}
			var controllerMap = [String: PlayerController]()
			for c in controllers {
				if c.online! {
					controllerMap[c.id] = c
				}
			}
			for btn in controllerButtons {
				guard let btnC = btn.controller, let c = controllerMap[btnC.id] else {
					btn.controller = nil
					continue
				}
				btn.controller = c
				controllerMap.removeValue(forKey: btnC.id)
			}
			for (_, v) in controllerMap {
				for btn in controllerButtons {
					if btn.controller == nil {
						btn.controller = v
						break
					}
				}
			}
			self.controllers = controllers
		} else if type == .NewMatch {
			HUD.hide()
			for btn in controllerButtons {
				btn.isSelected = false
			}
		} else if type == .Error {
			HUD.flash(.labeledError(title: json["msg"] as? String, subtitle: nil), delay: 1)
		}
	}
}

// MARK: swipe function
extension HallController: SWTableViewCellDelegate {
	fileprivate var rightButtons: [AnyObject] {
		let jumpButton = UIButton()
		jumpButton.setImage(UIImage(named: "CutLineButton"), for: UIControlState())
		jumpButton.backgroundColor = UIColor.clear
		let removeButton = UIButton()
		removeButton.setImage(UIImage(named: "RemoveTeamButton"), for: UIControlState())
		removeButton.backgroundColor = UIColor.clear
		return [jumpButton, removeButton]
	}

	func swipeableTableViewCell(_ cell: SWTableViewCell!, didTriggerRightUtilityButtonWith index: Int) {
		guard let team = teamFromCell(cell) else {
			return
		}
		if index == 0 {
			let json = JSON([
				"cmd": "teamCutLine",
				"teamID": team.id
			])
			WsClient.singleton.sendJSON(json)
		} else if index == 1 {
			let json = JSON([
				"cmd": "teamRemove",
				"teamID": team.id
			])
			WsClient.singleton.sendJSON(json)
		}
	}
	func swipeableTableViewCellShouldHideUtilityButtons(onSwipe cell: SWTableViewCell!) -> Bool {
		return true
	}
	func swipeableTableViewCell(_ cell: SWTableViewCell!, canSwipeTo state: SWCellState) -> Bool {
		guard let team = teamFromCell(cell) else {
			return false
		}
		if team.status != .waiting {
			return false
		}
		return true
	}

	fileprivate func teamFromCell(_ cell: SWTableViewCell) -> Team? {
		if let cellIndex = teamtableView.indexPath(for: cell), let tms = self.teams {
			return tms[cellIndex.row]
		}
		return nil
	}
}

extension HallController: UITableViewDataSource, UITableViewDelegate {
	func tableView(_ tableView: UITableView, numberOfRowsInSection section: Int) -> Int {
		return teams != nil ? teams!.count : 0
	}
	func numberOfSections(in tableView: UITableView) -> Int {
		return 1
	}
	func tableView(_ tableView: UITableView, cellForRowAt indexPath: IndexPath) -> UITableViewCell {
		let cell = tableView.dequeueReusableCell(withIdentifier: "HallTableViewCell")! as! HallTableViewCell
		let team = teams![indexPath.row]
		cell.setData(team, number: indexPath.row, active: team.id == topTeam?.id)
		cell.delegate = self
		cell.rightUtilityButtons = rightButtons
		return cell
	}
}
