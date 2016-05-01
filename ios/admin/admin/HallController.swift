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

class HallController: PLViewController {
	@IBOutlet weak var teamtableView: UITableView!
	@IBOutlet weak var teamIDLabel: UILabel!

	@IBOutlet weak var modeImageView: UIImageView!
	@IBOutlet weak var modeLabel: UILabel!
	@IBOutlet weak var playerNumberLabel: UILabel!
	@IBOutlet weak var readyButton: UIButton!
	@IBOutlet weak var startButton: UIButton!
	var refreshControl: UIRefreshControl!
	var teams: [Team]?
	var topTeam: Team?

	override func viewDidLoad() {
		super.viewDidLoad()
		refreshControl = UIRefreshControl()
		refreshControl.addTarget(self, action: #selector(HallController.refreshTeamData), forControlEvents: UIControlEvents.ValueChanged)
		teamtableView.addSubview(refreshControl)
	}
	override func viewWillAppear(animated: Bool) {
		super.viewWillAppear(animated)
		DataManager.singleton.subscriptData([.QueueData], receiver: self)
	}
	func refreshTeamData() {
		DataManager.singleton.queryData(.QueueData)
	}
	@IBAction func changeMode() {
		guard topTeam != nil else {
			return
		}
		let mode = topTeam!.mode == "g" ? "s" : "g"
		let json = JSON([
			"cmd": "teamChangeMode",
			"teamID": topTeam!.id,
			"mode": mode
		])
		WsClient.singleton.sendJSON(json)
	}
	@IBAction func callTeam(sender: UIButton) {
		guard topTeam != nil else {
			return
		}
		let json = JSON([
			"cmd": "teamCall",
			"teamID": topTeam!.id,
		])
		WsClient.singleton.sendJSON(json)
	}
	@IBAction func delayTeam(sender: UIButton) {
		guard topTeam != nil else {
			return
		}
		let json = JSON([
			"cmd": "teamDelay",
			"teamID": topTeam!.id,
		])
		WsClient.singleton.sendJSON(json)
	}
	@IBAction func addPlayer(sender: UIButton) {
		guard topTeam != nil && topTeam!.size < PLConstants.maxTeamSize else {
			return
		}
		let json = JSON([
			"cmd": "teamAddPlayer",
			"teamID": topTeam!.id,
		])
		WsClient.singleton.sendJSON(json)
	}
	@IBAction func removePlayer(sender: UIButton) {
		guard topTeam != nil && topTeam!.size > 1 else {
			return
		}
		let json = JSON([
			"cmd": "teamRemovePlayer",
			"teamID": topTeam!.id,
		])
		WsClient.singleton.sendJSON(json)
	}
	@IBAction func ready(sender: UIButton) {
		guard topTeam != nil else {
			return
		}
		let json = JSON([
			"cmd": "teamPrepare",
			"teamID": topTeam!.id,
		])
		WsClient.singleton.sendJSON(json)
	}
	@IBAction func start(sender: AnyObject) {
		guard topTeam != nil else {
			return
		}
		let json = JSON([
			"cmd": "teamStart",
			"teamID": topTeam!.id,
		])
		WsClient.singleton.sendJSON(json)
	}

	private func renderTopWaitingTeam() {
		guard topTeam != nil else {
			return
		}
		teamIDLabel.text = topTeam!.id
		playerNumberLabel.text = "\(topTeam!.size)"
		if topTeam!.mode == "g" {
			modeImageView.image = UIImage(named: "FunIcon")
			modeLabel.text = "[赏金模式]"
		} else {
			modeImageView.image = UIImage(named: "SurvivalIcon")
			modeLabel.text = "[生存模式]"
		}
		if topTeam!.status == .Waiting {
			readyButton.enabled = true
			startButton.enabled = false
		} else if topTeam!.status == .Prepare {
			readyButton.enabled = false
			startButton.enabled = true
		}
	}
}

extension HallController: DataReceiver {
	func onReceivedData(json: [String: AnyObject], type: DataType) {
		if type == .QueueData {
			teams = Mapper<Team>().mapArray(json["data"])
			if teams != nil {
				for team in teams! {
					if team.status == .Waiting || team.status == .Prepare {
						topTeam = team
						renderTopWaitingTeam()
						break
					}
				}
			}
			teamtableView.reloadData()
			refreshControl.endRefreshing()
		}
	}
}

// MARK: swipe function
extension HallController: SWTableViewCellDelegate {
	private var rightButtons: [AnyObject] {
		let jumpButton = UIButton()
		jumpButton.setImage(UIImage(named: "CutLineButton"), forState: .Normal)
		jumpButton.backgroundColor = UIColor.clearColor()
		let removeButton = UIButton()
		removeButton.setImage(UIImage(named: "RemoveTeamButton"), forState: .Normal)
		removeButton.backgroundColor = UIColor.clearColor()
		return [jumpButton, removeButton]
	}

	func swipeableTableViewCell(cell: SWTableViewCell!, didTriggerRightUtilityButtonWithIndex index: Int) {
		let team = teamFromCell(cell)
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
	func swipeableTableViewCellShouldHideUtilityButtonsOnSwipe(cell: SWTableViewCell!) -> Bool {
		return true
	}
	func swipeableTableViewCell(cell: SWTableViewCell!, canSwipeToState state: SWCellState) -> Bool {
		let team = teamFromCell(cell)
		if team.status != .Waiting {
			return false
		}
		return true
	}

	private func teamFromCell(cell: SWTableViewCell) -> Team {
		let cellIndex = teamtableView.indexPathForCell(cell)!
		return teams![cellIndex.row]
	}
}

extension HallController: UITableViewDataSource, UITableViewDelegate {
	func tableView(tableView: UITableView, numberOfRowsInSection section: Int) -> Int {
		return teams != nil ? teams!.count : 0
	}
	func numberOfSectionsInTableView(tableView: UITableView) -> Int {
		return 1
	}
	func tableView(tableView: UITableView, cellForRowAtIndexPath indexPath: NSIndexPath) -> UITableViewCell {
		let cell = tableView.dequeueReusableCellWithIdentifier("HallTableViewCell")! as! HallTableViewCell
		let team = teams![indexPath.row]
		cell.setData(team, number: indexPath.row, active: team.id == topTeam?.id)
		cell.delegate = self
		cell.rightUtilityButtons = rightButtons
		return cell
	}
}
