//
//  LoginViewController.swift
//  postgame
//
//  Created by tassar on 3/31/16.
//  Copyright © 2016 pulupulu. All rights reserved.
//

import UIKit
import Alamofire
import AlamofireObjectMapper
import AutoKeyboardScrollView
import SVProgressHUD
import EasyPeasy
import SwiftyUserDefaults
import PKHUD
// FIXME: comparison operators with optionals were removed from the Swift Standard Libary.
// Consider refactoring the code to use the non-optional operators.
fileprivate func < <T : Comparable>(lhs: T?, rhs: T?) -> Bool {
  switch (lhs, rhs) {
  case let (l?, r?):
    return l < r
  case (nil, _?):
    return true
  default:
    return false
  }
}

// FIXME: comparison operators with optionals were removed from the Swift Standard Libary.
// Consider refactoring the code to use the non-optional operators.
fileprivate func > <T : Comparable>(lhs: T?, rhs: T?) -> Bool {
  switch (lhs, rhs) {
  case let (l?, r?):
    return l > r
  default:
    return rhs < lhs
  }
}


let SegueIDShowMatchResult = "ShowMatchResult"

class LoginViewController: PLViewController {

	/*
	 为什么要这个wrapperView看下面
	 @link https://github.com/honghaoz/AutoKeyboardScrollView#work-with-interface-builder
	 */
	@IBOutlet weak var wrapperView: UIView!
	@IBOutlet weak var usernameTextField: LoginTextField!
	@IBOutlet weak var passwordTextField: LoginTextField!
	@IBOutlet weak var deviceIDTextField: LoginTextField!
	@IBOutlet weak var loginButton: UIButton!

	/**
	 双击登陆界面右上角出现配置窗口
	 */
	@IBAction func showConfig(_ sender: UITapGestureRecognizer) {
		let alert = UIAlertController(title: "设置", message: nil, preferredStyle: .alert)
		alert.addTextField { (textfield) in
			textfield.placeholder = Defaults[.host]
		}
		alert.addTextField { textfield in
			textfield.placeholder = Defaults[.deviceID]
		}
		alert.addTextField { textfield in
			textfield.placeholder = Defaults[.websiteHost]
		}
		let cancelAction = UIAlertAction(title: "取消", style: .cancel, handler: nil)
		alert.addAction(cancelAction)
		weak var weakAlert = alert
		let doneAction = UIAlertAction(title: "确定", style: .default) { (action) in
			if let host = weakAlert?.textFields![0].text, host != "" {
				Defaults[.host] = host
				WsClient.singleton.connect(PLConstants.getWsAddress())
			}
			if let num = weakAlert?.textFields![1].text, num != "" {
				Defaults[.deviceID] = num
			}
			if let host = weakAlert?.textFields![2].text, host != "" {
				Defaults[.websiteHost] = host
			}
		}
		alert.addAction(doneAction)
		present(alert, animated: true, completion: nil)
	}

	@IBAction func textFieldValueChanged(_ sender: UITextField) {
		if usernameTextField.text?.characters.count > 0 && passwordTextField.text?.characters.count > 0 && deviceIDTextField.text?.characters.count > 0 {
			self.loginButton.isEnabled = true
		} else {
			self.loginButton.isEnabled = false
		}
	}

	@IBAction func login() {
		HUD.show(.progress)
		let p = [
			"username": self.usernameTextField.text!,
			"password": self.passwordTextField.text!,
		]
        request(PLConstants.getWebsiteAddress("user/login"), method: .post, parameters: p, encoding: URLEncoding.default, headers: nil)
			.validate()
            .responseObject(completionHandler: { (resp: DataResponse<LoginResult>) in
				HUD.hide()
				if let _ = resp.result.error {
					HUD.flash(.error, delay: 2)
				} else {
					let m = resp.result.value!
					if m.code != nil && m.code == 0 {
						self.performSegue(withIdentifier: SegueIDShowMatchResult, sender: m)
					} else {
						HUD.flash(.labeledError(title: m.error, subtitle: nil), delay: 2)
					}
				}
		})
	}
	@IBAction func skip() {
		performSegue(withIdentifier: SegueIDShowMatchResult, sender: nil)
	}

	override func prepare(for segue: UIStoryboardSegue, sender: Any?) {
		if segue.identifier == SegueIDShowMatchResult {
			let vc = segue.destination as! MatchResultController
			vc.isAdmin = false
			vc.loginInfo = sender as? LoginResult
		}
	}
}

// MARK: UIViewController
extension LoginViewController {

	override func viewDidLoad() {
		super.viewDidLoad()
		let scrollView = AutoKeyboardScrollView()
		view.addSubview(scrollView)
		wrapperView.removeFromSuperview()
		scrollView.contentView.addSubview(wrapperView)
		scrollView.backgroundColor = wrapperView.backgroundColor
		scrollView.isUserInteractionEnabled = true
		scrollView.bounces = true
		scrollView.isScrollEnabled = true
		scrollView <- Edges()
		wrapperView <- Edges()
		scrollView.setTextMargin(175, forTextField: usernameTextField)
		scrollView.setTextMargin(140, forTextField: passwordTextField)
	}
}
