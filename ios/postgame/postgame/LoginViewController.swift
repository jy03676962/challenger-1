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

class LoginViewController: PLBaseViewController {

	var usernameTextField: UITextField?
	var passwordTextField: UITextField?
	var loginButton: UIButton?

	func formChanged(textField: UITextField) {
		if usernameTextField?.text?.characters.count > 0 && passwordTextField?.text?.characters.count > 0 {
			self.loginButton?.enabled = true
		} else {
			self.loginButton?.enabled = false
		}
	}

	func login() {
		let parameters: [String: AnyObject] = [
			"username": usernameTextField!.text!,
			"password": passwordTextField!.text!
		]
		Alamofire.request(.POST, "\(Constants.host)/login", parameters: parameters)
			.responseJSON { response in
				if let JSON = response.result.value {
					log.debug("\(JSON["username"]) has logined")
				}
		}
	}

	func skip() {
		log.debug("skip login")
	}
}

// MARK: UIViewController
extension LoginViewController {

	override func viewWillAppear(animated: Bool) {
		super.viewWillAppear(animated)
		let selector = #selector(LoginViewController.formChanged(_:))
		NSNotificationCenter.defaultCenter().addObserver(self, selector: selector, name: UITextFieldTextDidChangeNotification, object: nil)
	}

	override func viewDidDisappear(animated: Bool) {
		super.viewDidDisappear(animated)
		NSNotificationCenter.defaultCenter().removeObserver(self, name: UITextFieldTextDidChangeNotification, object: nil)
	}

	override func viewDidLoad() {
		super.viewDidLoad()
		setupViews()
	}

	private func setupViews() {
		let textFieldHeight: CGFloat = 35
		let xOffsetButton: CGFloat = 10
		let minMargin: CGFloat = 140.0
		let buttonSize = UIImage(named: "LoginButtonEnabled")!.size
		let scrollView = AutoKeyboardScrollView()
		let usernameTextField = UITextField()
		let passwordTextField = UITextField()
		let loginButton = UIButton()
		let skipButton = UIButton()
		scrollView.backgroundColor = UIColor.clearColor()
		view.addSubview(scrollView)
		scrollView.contentView.addSubview(usernameTextField)
		scrollView.contentView.addSubview(passwordTextField)
		scrollView.contentView.addSubview(loginButton)
		scrollView.contentView.addSubview(skipButton)
		func styleTextField(tf: UITextField, ph: String) -> () {
			tf.layer.borderColor = UIColor(rgba: "#4B6C87").CGColor
			tf.layer.borderWidth = 1
			tf.attributedPlaceholder = NSAttributedString(string: ph,
				attributes: [NSForegroundColorAttributeName: UIColor(rgba: "#424242"),])
			tf.leftView = UIView(frame: CGRect(x: 0, y: 0, width: 10, height: 0))
			tf.leftViewMode = .Always
			tf.clearButtonMode = .Always
			tf.textColor = UIColor.whiteColor()
			let btn = tf.valueForKey("_clearButton")
			btn?.setImage(UIImage(named: "TextFieldClear"), forState: .Normal)
		}
		styleTextField(usernameTextField, ph: "账号")
		styleTextField(passwordTextField, ph: "密码")
		passwordTextField.secureTextEntry = true
		loginButton.setBackgroundImage(UIImage(named: "LoginButtonEnabled"), forState: .Normal)
		loginButton.setBackgroundImage(UIImage(named: "LoginButtonDisabled"), forState: .Disabled)
		skipButton.setBackgroundImage(UIImage(named: "SkipButtonEnabled"), forState: .Normal)
		skipButton.setBackgroundImage(UIImage(named: "SkipButtonDisabled"), forState: .Disabled)
		loginButton.enabled = false

		// contraints
		scrollView.mas_makeConstraints({ m in
			m.edges.equalTo()(self.view)
		})
		func constraintTextField(tf: UITextField) -> () {
			tf.mas_makeConstraints({ m in
				m.centerX.equalTo()(scrollView)
				m.width.equalTo()(250)
				m.height.equalTo()(textFieldHeight)
			})
		}
		constraintTextField(usernameTextField)
		constraintTextField(passwordTextField)
		usernameTextField.mas_makeConstraints({ m in
			m.top.equalTo()(400)
		})
		scrollView.setTextMargin(textFieldHeight + minMargin, forTextField: usernameTextField)
		scrollView.setTextMargin(minMargin, forTextField: passwordTextField)
		passwordTextField.mas_makeConstraints({ m in
			m.top.equalTo()(usernameTextField.mas_bottom)
		})
		loginButton.mas_makeConstraints({ m in
			m.top.equalTo()(passwordTextField.mas_bottom).offset()(50)
			m.right.equalTo()(scrollView.mas_centerX).offset()(-10)
			m.width.equalTo()(buttonSize.width)
			m.height.equalTo()(buttonSize.height)
		})
		skipButton.mas_makeConstraints({ m in
			m.top.equalTo()(loginButton)
			m.left.equalTo()(scrollView.mas_centerX).offset()(10)
			m.width.equalTo()(buttonSize.width)
			m.height.equalTo()(buttonSize.height)
		})

		loginButton.addTarget(self, action: #selector(LoginViewController.login), forControlEvents: .TouchUpInside)
		skipButton.addTarget(self, action: #selector(LoginViewController.skip), forControlEvents: .TouchUpInside)

		self.usernameTextField = usernameTextField
		self.passwordTextField = passwordTextField
		self.loginButton = loginButton
	}
}
