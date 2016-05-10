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
import AlamofireObjectMapper
import SwiftyUserDefaults
import ObjectMapper

let SegueIDShowSurvey = "ShowSurvey"

class MatchResultController: PLViewController {
	var matchData: MatchData? {
		didSet {
			if self.viewIfLoaded != nil {
				renderData()
			}
		}
	}
	var playerData: PlayerData?
	var isAdmin: Bool = false

	@IBOutlet weak var headerImageView: UIImageView!
	@IBOutlet weak var tableHeaderImageView: UIImageView!

	@IBOutlet weak var playerTableView: UITableView!
	@IBOutlet weak var teamIDLabel: UILabel!
	@IBOutlet weak var scoreLabel: UILabel!

	@IBOutlet weak var stopAnswerButton: UIButton!
	@IBOutlet weak var startSurveyButton: UIButton!
	@IBOutlet weak var passSurveyButton: UIButton!
	@IBOutlet var playersLabel: [UILabel]!

	@IBAction func startSurvey() {
		HUD.show(.Progress)
		Alamofire.request(.GET, PLConstants.getHttpAddress("api/survey"))
			.validate()
			.responseObject(completionHandler: { (response: Response<Survey, NSError>) in
				HUD.hide()
				if let err = response.result.error {
					HUD.show(.LabeledError(title: err.localizedDescription, subtitle: nil))
				} else if let survey = response.result.value {
					self.performSegueWithIdentifier(SegueIDShowSurvey, sender: survey)
				}
		})
	}

	@IBAction func passSurvey() {
		let sb = UIStoryboard(name: "Main", bundle: nil)
		let login = sb.instantiateViewControllerWithIdentifier("LoginViewController")
		navigationController?.setViewControllers([login], animated: true)
	}

	@IBAction func endAnswer() {
		guard matchData != nil else {
			return
		}
		HUD.show(.Progress)
		Alamofire.request(.POST, PLConstants.getHttpAddress("api/stop_answer"), parameters: ["mid": matchData!.id], encoding: .URL, headers: nil)
			.validate()
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

	override func viewDidLoad() {
		super.viewDidLoad()
		self.playersLabel = self.playersLabel.sort({ (l1, l2) -> Bool in
			return l1.tag < l2.tag
		})
	}

	override func viewWillAppear(animated: Bool) {
		super.viewWillAppear(animated)
		adjustViews()
		if isAdmin {
			DataManager.singleton.subscribeData([.UpdatePlayerData], receiver: self)
			renderData()
		} else if matchData == nil {
			HUD.show(.Progress)
			Alamofire.request(.GET, PLConstants.getHttpAddress("api/answering"))
				.validate()
				.responseJSON(completionHandler: { response in
					HUD.hide()
					if let err = response.result.error {
						HUD.show(.LabeledError(title: err.localizedDescription, subtitle: nil))
					} else if let d = response.result.value {
						let code = d["code"] as! Int
						if code == 0 {
							self.matchData = Mapper<MatchData>().map(d["data"])
						}
						self.renderData()
					}
			})
		} else {
			renderData()
		}
	}

	override func viewDidDisappear(animated: Bool) {
		super.viewDidDisappear(animated)
		if isAdmin {
			DataManager.singleton.unsubscribe(self)
		}
	}

	override func prepareForSegue(segue: UIStoryboardSegue, sender: AnyObject?) {
		if segue.identifier == SegueIDShowSurvey {
			let vc = segue.destinationViewController as! SurveyController
			vc.playerData = self.playerData
			vc.survey = sender as! Survey
		}
	}

	func adjustViews() {
		if isAdmin {
			stopAnswerButton.hidden = false
			startSurveyButton.hidden = true
			passSurveyButton.hidden = true
			for label in self.playersLabel {
				label.hidden = false
			}
		} else {
			stopAnswerButton.hidden = true
			startSurveyButton.hidden = false
			passSurveyButton.hidden = false
			for label in self.playersLabel {
				label.hidden = true
			}
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
			if !isAdmin {
				for pd in data.member {
					if pd.cid.componentsSeparatedByString(":")[1] == Defaults[.deviceID] {
						self.playerData = pd
					}
				}
			} else {
				for (i, pd) in data.member.enumerate() {
					self.playersLabel[i].text = "\(pd.getName()): \(pd.answered)/\(Defaults[.qCount])"
				}
			}
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

extension MatchResultController: DataReceiver {
	func onReceivedData(json: [String: AnyObject], type: DataType) {
		if type == .UpdatePlayerData {
			if let data = matchData {
				let playerData = Mapper<PlayerData>().map(json["data"])!
				for (i, pd) in data.member.enumerate() {
					if pd.id == playerData.id {
						self.playersLabel[i].text = "\(playerData.getName()): \(playerData.answered)/\(Defaults[.qCount])"
						data.member[i] = playerData
						break
					}
				}
			}
		}
	}
}
