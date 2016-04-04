//
//  LoginViewController.swift
//  postgame
//
//  Created by tassar on 3/31/16.
//  Copyright © 2016 pulupulu. All rights reserved.
//

import UIKit
import AutoKeyboardScrollView

class LoginViewController: PLBaseViewController {

	var usernameTextField: UITextField?
	var passwordTextField: UITextField?
	var loginButton: UIButton?

	override func backgroundImage() -> UIImage? {
		return UIImage(named: "GlobalBackground")
	}

	func textFiledChanged(textField: UITextField) {
		if usernameTextField?.text?.characters.count > 0 && passwordTextField?.text?.characters.count > 0 {
			self.loginButton?.enabled = true
		} else {
			self.loginButton?.enabled = false
		}
	}

	override func viewDidLoad() {
		super.viewDidLoad()
		setupViews()
	}

	override func viewWillAppear(animated: Bool) {
		super.viewWillAppear(animated)
		navigationController?.navigationBar.hidden = true
		let selector = #selector(LoginViewController.textFiledChanged)
		NSNotificationCenter.defaultCenter().addObserver(self, selector: selector, name: UITextFieldTextDidChangeNotification, object: nil)
	}

	override func viewDidDisappear(animated: Bool) {
		super.viewDidDisappear(animated)
		NSNotificationCenter.defaultCenter().removeObserver(self, name: UITextFieldTextDidChangeNotification, object: nil)
	}

	override func prefersStatusBarHidden() -> Bool {
		return true
	}

	private func setupViews() {
		let textFieldHeight: CGFloat = 35
		let minMargin: CGFloat = 140.0
		let buttonSize = UIImage(named: "LoginButtonEnabled")!.size
		let scrollView = AutoKeyboardScrollView()
		let usernameTextField = UITextField()
		let passwordTextField = UITextField()
		let buttonContainer = UIView()
		let loginButton = UIButton()
		let skipButton = UIButton()
		scrollView.backgroundColor = UIColor.clearColor()
		view.addSubview(scrollView)
		scrollView.contentView.addSubview(usernameTextField)
		scrollView.contentView.addSubview(passwordTextField)
		scrollView.contentView.addSubview(buttonContainer)
		buttonContainer.addSubview(loginButton)
		buttonContainer.addSubview(skipButton)
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
		loginButton.setImage(UIImage(named: "LoginButtonEnabled"), forState: .Normal)
		loginButton.setImage(UIImage(named: "LoginButtonDisabled"), forState: .Disabled)
		skipButton.setImage(UIImage(named: "SkipButtonEnabled"), forState: .Normal)
		skipButton.setImage(UIImage(named: "SkipButtonDisabled"), forState: .Disabled)
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
		buttonContainer.mas_makeConstraints({ m in
			m.centerX.equalTo()(scrollView)
			m.top.equalTo()(passwordTextField.mas_bottom).with().offset()(50)
		})
		loginButton.mas_makeConstraints({ m in
			m.top.equalTo()(buttonContainer)
			m.left.equalTo()(buttonContainer)
			m.width.equalTo()(buttonSize.width)
			m.height.equalTo()(buttonSize.height)
		})
		skipButton.mas_makeConstraints({ m in
			m.left.equalTo()(loginButton.mas_right).with().offset()(15)
			m.top.equalTo()(buttonContainer)
			m.right.equalTo()(buttonContainer)
			m.width.equalTo()(buttonSize.width)
			m.height.equalTo()(buttonSize.height)
		})

		self.usernameTextField = usernameTextField
		self.passwordTextField = passwordTextField
		self.loginButton = loginButton
	}
}
