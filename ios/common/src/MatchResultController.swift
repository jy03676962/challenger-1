//
//  MatchResultController.swift
//  admin
//
//  Created by tassar on 5/6/16.
//  Copyright © 2016 pulupulu. All rights reserved.
//

import UIKit
import PKHUD
import Alamofire

class MatchResultController: PLViewController {
	var matchData: MatchData? {
		didSet {
			if self.viewIfLoaded != nil {
				renderData()
			}
		}
	}
	var isAdmin: Bool = false

	@IBOutlet weak var headerImageView: UIImageView!
	@IBOutlet weak var tableHeaderImageView: UIImageView!

	@IBOutlet weak var playerTableView: UITableView!
	@IBOutlet weak var teamIDLabel: UILabel!
	@IBOutlet weak var scoreLabel: UILabel!

	@IBOutlet weak var stopAnswerButton: UIButton!
	@IBAction func endAnswer() {
		guard matchData != nil else {
			return
		}
		HUD.show(.Progress)
		Alamofire.request(.POST, PLConstants.getHttpAddress("api/stop_answer"), parameters: ["mid": matchData!.id], encoding: .URL, headers: nil)
			.responseJSON(completionHandler: { res in
				HUD.hide()
				if let err = res.result.error {
					HUD.flash(.LabeledError(title: err.localizedDescription, subtitle: nil), delay: 2)
				} else if let d = res.result.value as? UInt {
					if d == self.matchData?.id {
						self.presentingViewController?.dismissViewControllerAnimated(true, completion: nil)
					}
				}
		})
	}

	override func viewWillAppear(animated: Bool) {
		super.viewWillAppear(animated)
		adjust()
		renderData()
	}

	func adjust() {
		if isAdmin {
			stopAnswerButton.hidden = false
		} else {
			stopAnswerButton.hidden = true
		}
	}

	func renderData() {
		if let data = matchData {
			HUD.hide()
			if data.mode == "g" {
				headerImageView.image = UIImage(named: "FunImage")
				tableHeaderImageView.image = UIImage(named: "MatchGoldResultHeader")
			} else {
				headerImageView.image = UIImage(named: "SurvivalImage")
				tableHeaderImageView.image = UIImage(named: "MatchResultHeader")
			}
			teamIDLabel.text = data.teamID
			playerTableView.reloadData()
		} else {
			HUD.show(.LabeledProgress(title: "等待数据中...", subtitle: nil))
		}
	}
}

extension MatchResultController: UITableViewDataSource, UITableViewDelegate {
	func numberOfSectionsInTableView(tableView: UITableView) -> Int {
		return 1
	}

	func tableView(tableView: UITableView, numberOfRowsInSection section: Int) -> Int {
		return matchData == nil ? 0 : matchData!.member.count
	}

	func tableView(tableView: UITableView, cellForRowAtIndexPath indexPath: NSIndexPath) -> UITableViewCell {
		let cell = tableView.dequeueReusableCellWithIdentifier("MatchResultCell") as! MatchResultCell
		let data = matchData!.member[indexPath.row]
		cell.setData(data)
		return cell
	}
}
