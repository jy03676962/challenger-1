//
//  LoginViewController.swift
//  postgame
//
//  Created by tassar on 3/31/16.
//  Copyright © 2016 pulupulu. All rights reserved.
//

import UIKit
import Alamofire
import AutoKeyboardScrollView
import SVProgressHUD
import EasyPeasy
import SwiftyUserDefaults

let SegueIDShowMatchResult = "ShowMatchResult"

class LoginViewController: PLViewController {

	/*
	 为什么要这个wrapperView看下面
	 @link https://github.com/honghaoz/AutoKeyboardScrollView#work-with-interface-builder
	 */
	@IBOutlet weak var wrapperView: UIView!
	@IBOutlet weak var usernameTextField: LoginTextField!
	@IBOutlet weak var passwordTextField: LoginTextField!
	@IBOutlet weak var loginButton: UIButton!

	/**
	 双击登陆界面右上角出现配置窗口
	 */
	@IBAction func showConfig(sender: UITapGestureRecognizer) {
		let alert = UIAlertController(title: "设置", message: nil, preferredStyle: .Alert)
		alert.addTextFieldWithConfigurationHandler { (textfield) in
			textfield.placeholder = "输入HOST"
		}
		alert.addTextFieldWithConfigurationHandler { textfield in
			textfield.placeholder = "输入编号"
		}
		let cancelAction = UIAlertAction(title: "取消", style: .Cancel, handler: nil)
		alert.addAction(cancelAction)
		weak var weakAlert = alert
		let doneAction = UIAlertAction(title: "确定", style: .Default) { (action) in
			if let host = weakAlert?.textFields![0].text {
				Defaults[.host] = host
				WsClient.singleton.connect(PLConstants.getWsAddress())
			}
			if let num = weakAlert?.textFields![1].text {
				Defaults[.deviceID] = num
			}
		}
		alert.addAction(doneAction)
		presentViewController(alert, animated: true, completion: nil)
	}

	@IBAction func usernameEditEnd(sender: UITextField) {
		passwordTextField.becomeFirstResponder()
	}

	@IBAction func passwordEditEnd(sender: UITextField) {
		login()
	}

	@IBAction func textFieldValueChanged(sender: UITextField) {
		if usernameTextField.text?.characters.count > 0 && passwordTextField.text?.characters.count > 0 {
			self.loginButton.enabled = true
		} else {
			self.loginButton.enabled = false
		}
	}

	@IBAction func login() {
//		let parameters: [String: AnyObject] = [
//			"username": usernameTextField.text!,
//			"password": passwordTextField.text!
//		]
//		SVProgressHUD.show()
		// TODO login
//		Alamofire.request(.POST, PLConstants.getHttpAddress("api/login"), parameters: parameters)
//			.responseJSON { response in
//				SVProgressHUD.dismiss()
//				if let JSON = response.result.value {
//					log.debug("\(JSON["username"]) has logined")
//				}
//		}
	}
	@IBAction func skip() {
		performSegueWithIdentifier(SegueIDShowMatchResult, sender: nil)
	}

	override func prepareForSegue(segue: UIStoryboardSegue, sender: AnyObject?) {
		if segue.identifier == SegueIDShowMatchResult {
			let vc = segue.destinationViewController as! MatchResultController
			let app = UIApplication.sharedApplication().delegate as! AppDelegate
			vc.matchData = app.matchData
			vc.isAdmin = false
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
		scrollView.userInteractionEnabled = true
		scrollView.bounces = true
		scrollView.scrollEnabled = true
		scrollView <- Edges()
		wrapperView <- Edges()
		scrollView.setTextMargin(175, forTextField: usernameTextField)
		scrollView.setTextMargin(140, forTextField: passwordTextField)
	}
}
