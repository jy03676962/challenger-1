//
//  HistoryController.swift
//  admin
//
//  Created by tassar on 4/20/16.
//  Copyright © 2016 pulupulu. All rights reserved.
//

import UIKit
import Alamofire
import AlamofireObjectMapper
import PKHUD
import SwiftyJSON

let segueIDPresentMatchResult = "PresentMatchResult"

class HistoryController: PLViewController {
	@IBOutlet weak var tableView: UITableView!
	@IBOutlet weak var startAnswerButton: UIButton!

	var refreshControl: UIRefreshControl!
	var data: [MatchData]?
	var isAnswering: Int? {
		guard let d = data else {
			return nil
		}
		for (i, m) in d.enumerated() {
			if m.answerType == .answering {
				return i
			}
		}
		return nil
	}

	@IBAction func startAnswer() {
		guard let data = self.data, let indexPaths = tableView.indexPathsForSelectedRows, indexPaths.count == 1 else {
			return
		}
		if let ing = isAnswering {
			if ing != indexPaths[0].row {
				HUD.flash(.labeledError(title: "有其他正在答题的队伍", subtitle: "请先结束该组答题后重试"), delay: 1)
				return
			}
		}
		let matchData = data[indexPaths[0].row]
		if matchData.eid != nil && matchData.eid != "" {
			self.startAnswerAfterAdd(matchData)
		} else {
			HUD.show(.progress)
			var playerDataList: [Any] = []
			for player in matchData.member {
				let pd: [String: Any] = [
					"player_id": player.cid,
					"player_score": String(player.gold - player.lostGold),
					"player_catch": String(player.hitCount),
					"player_rank": player.grade,
				]
				playerDataList.append(pd)
			}
			let p: [String: Any] = [
				"mode": matchData.mode == "g" ? "0" : "1",
				"time": String(Int(matchData.elasped * 1000)),
				"gold": String(matchData.gold),
				"player_num": String(matchData.member.count),
				"team_rampage": String(matchData.rampageCount),
				"team_rank": matchData.grade,
				"player_data": JSON(playerDataList).rawString()!,
			]
            request(PLConstants.getWebsiteAddress("challenger/match"), method: .post, parameters: p, encoding: URLEncoding.default, headers: nil)
				.validate()
                .responseObject(completionHandler: { (response: DataResponse<AddMatchResult>) in
					HUD.hide()
					if let err = response.result.error {
						HUD.flash(.labeledError(title: err.localizedDescription, subtitle: nil), delay: 2)
					} else {
						if response.result.value?.code != 0 {
							HUD.flash(.labeledError(title: response.result.value?.error, subtitle: nil), delay: 2)
						} else {
							matchData.eid = String(response.result.value!.matchID)
							self.startAnswerAfterAdd(matchData)
						}
					}
			})
		}
	}

	override func viewDidLoad() {
		super.viewDidLoad()
		refreshControl = UIRefreshControl()
		refreshControl.addTarget(self, action: #selector(refreshHistory), for: .valueChanged)
		tableView.addSubview(refreshControl)
	}

	override func viewWillAppear(_ animated: Bool) {
		super.viewWillAppear(animated)
		refreshHistory()
	}

	func refreshHistory() {
        request(PLConstants.getHttpAddress("api/history"))
			.validate()
            .responseArray(completionHandler: { (response: DataResponse<[MatchData]>) in
				self.data = response.result.value
				if self.data != nil {
					self.tableView.reloadData()
				}
				self.refreshControl.endRefreshing()
		})
	}

	func startAnswerAfterAdd(_ matchData: MatchData) {
		HUD.show(.progress)
        let p: [String: Any] = [
            "mid": matchData.id,
            "eid": matchData.eid!
        ]
        request(PLConstants.getHttpAddress("api/start_answer"), method: .post, parameters: p, encoding: URLEncoding.default, headers: nil)
			.validate()
            .responseObject(completionHandler: { (response: DataResponse<MatchData>) in
				HUD.hide()
				if let err = response.result.error {
					HUD.flash(.labeledError(title: err.localizedDescription, subtitle: nil), delay: 2)
				} else {
					self.performSegue(withIdentifier: segueIDPresentMatchResult, sender: matchData)
				}
		})
	}

	override func prepare(for segue: UIStoryboardSegue, sender: Any?) {
		if segue.identifier == segueIDPresentMatchResult {
			let vc = segue.destination as! MatchResultController
			vc.isAdmin = true
			vc.matchData = sender as? MatchData
		}
	}
}

extension HistoryController: UITableViewDelegate, UITableViewDataSource {
	func tableView(_ tableView: UITableView, numberOfRowsInSection section: Int) -> Int {
		return self.data == nil ? 0 : self.data!.count
	}
	func numberOfSections(in tableView: UITableView) -> Int {
		return 1
	}
	func tableView(_ tableView: UITableView, cellForRowAt indexPath: IndexPath) -> UITableViewCell {
		let cell = tableView.dequeueReusableCell(withIdentifier: "HistoryTableViewCell") as! HistoryTableViewCell
		cell.setData(data![indexPath.row])
		return cell
	}
	func tableView(_ tableView: UITableView, didSelectRowAt indexPath: IndexPath) {
		startAnswerButton.isEnabled = true
	}
	func tableView(_ tableView: UITableView, didDeselectRowAt indexPath: IndexPath) {
		startAnswerButton.isEnabled = false
	}
}
