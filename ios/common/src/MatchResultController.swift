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
			guard !isAdmin else {
				return
			}
			if let data = matchData {
				let idx: Int = Int(Defaults[.deviceID])!
				self.playerData = data.member[idx - 1]
			}
		}
	}
	var playerData: PlayerData?
	var loginInfo: LoginResult?
	var isAdmin: Bool = false
	var surveyStarted: Bool = false
	var showAnswerStatus: Bool {
		return isAdmin
	}

	@IBOutlet weak var headerImageView: UIImageView!
	@IBOutlet weak var tableHeaderImageView: UIImageView!

	@IBOutlet weak var playerTableView: UITableView!
	@IBOutlet weak var teamIDLabel: UILabel!
	@IBOutlet weak var scoreLabel: UILabel!
	@IBOutlet weak var personalScoreLabel: UILabel!

	@IBOutlet weak var stopAnswerButton: UIButton!
	@IBOutlet weak var startSurveyButton: UIButton!
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
					self.surveyStarted = true
				}
		})
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
		DataManager.singleton.subscribeData([.UpdatePlayerData], receiver: self)
		if isAdmin {
			renderData()
		} else if matchData == nil {
			HUD.show(.Progress)
			Alamofire.request(.GET, PLConstants.getHttpAddress("api/answering"))
				.validate()
				.responseJSON(completionHandler: { response in
					HUD.hide()
					if let _ = response.result.error {
						self.waitingData()
					} else if let d = response.result.value {
						let code = d["code"] as! Int
						if code == 0 {
							self.matchData = Mapper<MatchData>().map(d["data"])
							self.uploadAndRender()
						} else {
							self.waitingData()
						}
					}
			})
		} else {
			uploadAndRender()
		}
	}

	override func viewDidDisappear(animated: Bool) {
		super.viewDidDisappear(animated)
		DataManager.singleton.unsubscribe(self)
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
		} else {
			stopAnswerButton.hidden = true
			startSurveyButton.hidden = self.surveyStarted
			for label in self.playersLabel {
				label.hidden = true
			}
		}
	}

	func uploadAndRender() {
		if let userInfo = self.loginInfo {
			let p: [String: AnyObject] = [
				"match_id": self.matchData!.id,
				"user_id": userInfo.userID,
				"username": userInfo.username,
			]
			Alamofire.request(.POST, PLConstants.getWebsiteAddress("challenger/adduser"), parameters: p, encoding: .URL, headers: nil)
				.validate()
				.responseObject(completionHandler: { (resp: Response<BaseResult, NSError>) in
					log.debug(resp.debugDescription)
			})
			if let pd = playerData {
				Alamofire.request(.POST,
					PLConstants.getHttpAddress("api/update_player"),
					parameters: ["pid": String(pd.id), "name": userInfo.username, "eid": String(userInfo.userID)],
					encoding: .URL,
					headers: nil)
					.validate()
					.responseData(completionHandler: { resp in
						log.debug(resp.debugDescription)
				})
			}
		}
		renderData()
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
			if isAdmin {
				for (i, label) in self.playersLabel.enumerate() {
					if i < data.member.count {
						let pd = data.member[i]
						label.text = "\(pd.getName()): \(pd.answered) / \(Defaults[.qCount]) "
						label.hidden = false
					} else {
						label.hidden = true
					}
				}
			}
		}
	}

	func waitingData() {
		DataManager.singleton.subscribeData([.StartAnswer], receiver: self)
		HUD.show(.LabeledProgress(title: "等待数据中 ... ", subtitle: nil))
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
				if showAnswerStatus {
					for (i, pd) in data.member.enumerate() {
						if pd.id == playerData.id {
							self.playersLabel[i].text = "\(playerData.getName()): \(playerData.answered) / \(Defaults[.qCount]) "
							data.member[i] = playerData
							break
						}
					}
				}
				self.playerTableView.reloadData()
			}
		} else if type == .StartAnswer {
			HUD.hide()
			matchData = Mapper<MatchData>().map(json["data"])
			uploadAndRender()
		}
	}
}
