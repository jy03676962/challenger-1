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

class HallController: PLViewController {
	@IBOutlet weak var teamtableView: UITableView!
	@IBOutlet weak var teamIDLabel: UILabel!

	@IBOutlet weak var modeImageView: UIImageView!
	@IBOutlet weak var modeLabel: UILabel!
	@IBOutlet weak var playerNumberLabel: UILabel!
	@IBOutlet weak var readyButton: UIButton!
	@IBOutlet weak var startButton: UIButton!
	var refreshControl: UIRefreshControl!
	var teams: [JSON]?

	override func viewDidLoad() {
		super.viewDidLoad()
		refreshControl = UIRefreshControl()
		refreshControl.addTarget(self, action: #selector(HallController.refreshTeamData), forControlEvents: UIControlEvents.ValueChanged)
		teamtableView.addSubview(refreshControl)
	}
	override func viewWillAppear(animated: Bool) {
		super.viewWillAppear(animated)
		DataManager.singleton.subscriptData([.HallData], receiver: self)
	}
	func refreshTeamData() {
		DataManager.singleton.queryData(.HallData)
	}
	@IBAction func changeMode(sender: UITapGestureRecognizer) {
	}
	@IBAction func callTeam(sender: UIButton) {
	}
	@IBAction func delayTeam(sender: UIButton) {
	}
	@IBAction func addPlayer(sender: UIButton) {
	}
	@IBAction func removePlayer(sender: UIButton) {
	}
	@IBAction func ready(sender: UIButton) {
	}
	@IBAction func start(sender: AnyObject) {
	}
}

extension HallController: DataReceiver {
	func onReceivedData(json: JSON, type: DataType) {
		if type == .HallData {
			teams = json["teams"].array
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
		let teamJson = teamJsonFromCell(cell)
		if index == 0 {
			let json = JSON([
				"cmd": "teamCutLine",
				"teamID": teamJson["id"].stringValue
			])
			WsClient.singleton.sendJSON(json)
		} else if index == 1 {
			let json = JSON([
				"cmd": "teamRemove",
				"teamID": teamJson["id"].stringValue
			])
			WsClient.singleton.sendJSON(json)
		}
	}
	func swipeableTableViewCellShouldHideUtilityButtonsOnSwipe(cell: SWTableViewCell!) -> Bool {
		return true
	}
	func swipeableTableViewCell(cell: SWTableViewCell!, canSwipeToState state: SWCellState) -> Bool {
		let teamJson = teamJsonFromCell(cell)
		let status = TeamStatus(rawValue: teamJson["status"].intValue)!
		if status != .Waiting {
			return false
		}
		return true
	}

	private func teamJsonFromCell(cell: SWTableViewCell) -> JSON {
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
		cell.setData(teams![indexPath.row])
		cell.delegate = self
		cell.rightUtilityButtons = rightButtons
		return cell
	}
}
