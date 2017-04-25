//
//  QuestionTableViewController.swift
//  postgame
//
//  Created by tassar on 5/8/16.
//  Copyright © 2016 pulupulu. All rights reserved.
//

import UIKit

protocol QuestionViewControllerDelegate: class {
	func okAction(_ sender: QuestionViewController, answer: String)
}

class QuestionViewController: UIViewController {
	@IBOutlet weak var titleLabel: UILabel!
	@IBOutlet weak var okButton: UIButton!
	@IBOutlet weak var tableView: UITableView!
	var question: SurveyQuestion!
	var questionIndex: Int = 0
	var isLastQuestion = false
	weak var delegate: QuestionViewControllerDelegate?

	@IBAction func okAction() {
		let idx = tableView.indexPathsForSelectedRows![0].row
		var ans = ""
		ans.append(Character(UnicodeScalar(idx + 17)!))
		delegate?.okAction(self, answer: ans)
	}

	override func viewDidLoad() {
		super.viewDidLoad()
		titleLabel.text = question.q
		tableView.reloadData()
		if isLastQuestion {
			okButton.setBackgroundImage(UIImage(named: "SurveyDone"), for: UIControlState())
		} else {
			okButton.setBackgroundImage(UIImage(named: "SurveyOK"), for: UIControlState())
		}
	}
}

extension QuestionViewController: UITableViewDataSource, UITableViewDelegate {

	func numberOfSections(in tableView: UITableView) -> Int {
		return 1
	}

	func tableView(_ tableView: UITableView, numberOfRowsInSection section: Int) -> Int {
		return question.options.count
	}
	func tableView(_ tableView: UITableView, cellForRowAt indexPath: IndexPath) -> UITableViewCell {
		let cell = tableView.dequeueReusableCell(withIdentifier: "QuestionTableViewCell") as! QuestionTableViewCell
		cell.setData(question.options[indexPath.row])
		return cell
	}
	func tableView(_ tableView: UITableView, didSelectRowAt indexPath: IndexPath) {
		okButton.isEnabled = true
	}
	func tableView(_ tableView: UITableView, didDeselectRowAt indexPath: IndexPath) {
		okButton.isEnabled = false
	}
}
