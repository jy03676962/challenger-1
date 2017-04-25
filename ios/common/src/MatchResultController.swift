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
import QRCode

let SegueIDShowSurvey = "ShowSurvey"

class MatchResultController: PLViewController {
	var matchData: MatchData?
	var playerData: PlayerData? {
		guard let data = matchData, !isAdmin else {
			return nil
		}
		let idx: Int = Int(Defaults[.deviceID])!
		return idx <= data.member.count ? data.member[idx - 1]: nil
	}
	var loginInfo: LoginResult?
	var isAdmin: Bool = false
	var surveyStarted: Bool = false
	var showAnswerStatus: Bool {
		return isAdmin
	}

	@IBOutlet weak var headerImageView: UIImageView!
	@IBOutlet weak var tableHeaderImageView: UIImageView!
	@IBOutlet weak var QRCodeImageView: UIImageView!

	@IBOutlet weak var playerTableView: UITableView!
	@IBOutlet weak var teamIDLabel: UILabel!
	@IBOutlet weak var scoreLabel: UILabel!
	@IBOutlet weak var personalScoreLabel: UILabel!
	@IBOutlet weak var personalScoreHeader: UILabel!

	@IBOutlet weak var stopAnswerButton: UIButton!
	@IBOutlet weak var startSurveyButton: UIButton!
	@IBOutlet var playersLabel: [UILabel]!

	@IBAction func startSurvey() {
		HUD.show(.progress)
        Alamofire.request(PLConstants.getHttpAddress("api/survey"), method: .get)
			.validate()
			.responseObject(completionHandler: { (response: DataResponse<Survey>) in
				HUD.hide()
				if let err = response.result.error {
					HUD.show(.labeledError(title: err.localizedDescription, subtitle: nil))
				} else if let survey = response.result.value {
					self.performSegue(withIdentifier: SegueIDShowSurvey, sender: survey)
					self.surveyStarted = true
				}
		})
	}

	@IBAction func endAnswer() {
		guard matchData != nil else {
			return
		}
		HUD.show(.progress)
        let p : [String : Any] = [
            "mid" : matchData!.id
        ]
        request(PLConstants.getHttpAddress("api/stop_answer"), method: .post, parameters: p, encoding: URLEncoding.default, headers: nil)
            .validate()
            .responseJSON(completionHandler: { res in
                HUD.hide()
                if let err = res.result.error {
                    HUD.flash(.labeledError(title: err.localizedDescription, subtitle: nil), delay: 2)
                } else if let d = res.result.value as? UInt {
                    if d == self.matchData?.id {
                        self.navigationController?.popViewController(animated: true)
                    }
                }
            })
	}

	override func viewDidLoad() {
		super.viewDidLoad()
		self.playersLabel = self.playersLabel.sorted(by: { (l1, l2) -> Bool in
			return l1.tag < l2.tag
		})
	}

	override func viewWillAppear(_ animated: Bool) {
		super.viewWillAppear(animated)
		self.navigationItem.hidesBackButton = true
		adjustViews()
		DataManager.singleton.subscribeData([.UpdatePlayerData], receiver: self)
		if isAdmin {
			renderData()
		} else if matchData == nil {
			HUD.show(.progress)
			Alamofire.request(PLConstants.getHttpAddress("api/answering"))
				.validate()
                .responseObject(completionHandler: { (res: DataResponse<AnsweringResponse>) in
                    HUD.hide()
                    if let _ = res.result.error {
						self.waitingData()
                    } else if let answeringResponse = res.result.value {
                        if answeringResponse.code == 0 {
                            self.matchData = answeringResponse.data
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

	override func viewDidDisappear(_ animated: Bool) {
		super.viewDidDisappear(animated)
		DataManager.singleton.unsubscribe(self)
	}

	override func prepare(for segue: UIStoryboardSegue, sender: Any?) {
		if segue.identifier == SegueIDShowSurvey {
			let vc = segue.destination as! SurveyController
			vc.playerData = self.playerData
			vc.survey = sender as! Survey
		}
	}

	func adjustViews() {
		if isAdmin {
			stopAnswerButton.isHidden = false
			startSurveyButton.isHidden = true
		} else {
			stopAnswerButton.isHidden = true
			startSurveyButton.isHidden = self.surveyStarted
			for label in self.playersLabel {
				label.isHidden = true
			}
		}
	}

	func uploadAndRender() {
		if let userInfo = self.loginInfo, let pd = playerData {
			let p: [String: Any] = [
				"match_id": self.matchData!.eid!,
				"user_id": userInfo.userID,
				"username": userInfo.username,
				"player_id": pd.cid,
			]
            request(PLConstants.getWebsiteAddress("challenger/adduser"), method: .post, parameters: p, encoding: URLEncoding.default, headers: nil)
                .validate()
                .responseObject(completionHandler: { (resp: DataResponse<AddUserResult>) in
                    log.debug(resp.debugDescription)
                    guard let data = self.matchData else {
                        return
                    }
                    for p in data.member {
                        if p.id == pd.id {
                            p.level = resp.result.value?.level
                            p.url = resp.result.value?.url
                            break
                        }
                    }
                    self.renderData()
                })
			request(
				PLConstants.getHttpAddress("api/update_player"),
				method: .post,
				parameters: ["pid": String(pd.id), "name": userInfo.username, "eid": String(userInfo.userID)],
				encoding: URLEncoding.default,
				headers: nil)
				.validate()
				.responseData(completionHandler: { resp in
					log.debug(resp.debugDescription)
			})
		}
		renderData()
	}

	func renderData() {
		if let data = matchData {
			HUD.hide()
			if data.mode == "g" {
				headerImageView.image = UIImage(named: "FunImage")
				tableHeaderImageView.image = UIImage(named: "MatchGoldResultHeader")
				scoreLabel.text = "\(data.gold)G"
			} else {
				headerImageView.image = UIImage(named: "SurvivalImage")
				tableHeaderImageView.image = UIImage(named: "MatchResultHeader")
				scoreLabel.text = String(format: "%.2fS", data.elasped)
			}
			teamIDLabel.text = data.teamID
			playerTableView.reloadData()
			if isAdmin {
				for (i, label) in self.playersLabel.enumerated() {
					if i < data.member.count {
						let pd = data.member[i]
						label.text = "\(pd.getName()): \(pd.answered) / \(Defaults[.qCount]) "
						label.isHidden = false
					} else {
						label.isHidden = true
					}
				}
				self.personalScoreLabel.isHidden = true
				self.personalScoreHeader.isHidden = true
			} else if let pd = self.playerData {
				self.personalScoreLabel.isHidden = false
				self.personalScoreHeader.isHidden = false
				self.personalScoreLabel.text = "\(pd.gold - pd.lostGold)G"
				guard let url = pd.url else {
					return
				}
				var qrCode = QRCode(url)
				qrCode?.size = CGSize(width: 90, height: 90)
				self.QRCodeImageView.image = qrCode?.image
			}
		}
	}

	func waitingData() {
		DataManager.singleton.subscribeData([.StartAnswer], receiver: self)
		HUD.show(.labeledProgress(title: "等待数据中 ... ", subtitle: nil))
	}
}

extension MatchResultController: UITableViewDataSource, UITableViewDelegate {
	func numberOfSections(in tableView: UITableView) -> Int {
		return 1
	}

	func tableView(_ tableView: UITableView, numberOfRowsInSection section: Int) -> Int {
		return 4
	}

	func tableView(_ tableView: UITableView, cellForRowAt indexPath: IndexPath) -> UITableViewCell {
		let cell = tableView.dequeueReusableCell(withIdentifier: "MatchResultCell") as! MatchResultCell
		if matchData == nil {
			cell.setData(nil, current: false)
		} else if indexPath.row < matchData!.member.count {
			let data = matchData!.member[indexPath.row]
			var current = false
			if let pd = self.playerData {
				current = data.cid == pd.cid
			}
			cell.setData(data, current: current)
		} else {
			cell.setData(nil, current: false)
		}
		return cell
	}
}

extension MatchResultController: DataReceiver {
	func onReceivedData(_ json: [String: Any], type: DataType) {
		if type == .UpdatePlayerData {
			if let data = matchData {
                let playerData = Mapper<PlayerData>().map(JSONObject:json["data"])!
				for (i, pd) in data.member.enumerated() {
					if pd.id == playerData.id {
						if showAnswerStatus {
							self.playersLabel[i].text = "\(playerData.getName()): \(playerData.answered) / \(Defaults[.qCount]) "
						}
						data.member[i] = playerData
						break
					}
				}
				self.playerTableView.reloadData()
			}
		} else if type == .StartAnswer {
			DataManager.singleton.unsubscribe(self, type: .StartAnswer)
			HUD.hide()
            matchData = Mapper<MatchData>().map(JSONObject:json["data"])
			uploadAndRender()
		}
	}
}
