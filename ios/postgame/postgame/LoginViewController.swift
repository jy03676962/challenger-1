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
    
    private let scrollView = AutoKeyboardScrollView()
    private let usernameTextField = UITextField()
    private let passwordTextField = UITextField()
    
    override func backgroundImage() -> UIImage? {
        return UIImage(named: "GlobalBackground")
    }
}

// MARK: viewcontroller methods
extension LoginViewController {
    override func viewDidLoad() {
        super.viewDidLoad()
        setupViews()
        setupConstraints()
    }
    
    override func viewWillAppear(animated: Bool) {
        super.viewWillAppear(animated)
        navigationController?.navigationBar.hidden = true
    }
    
    override func prefersStatusBarHidden() -> Bool {
        return true
    }
    
    private func setupViews() {
        scrollView.backgroundColor = UIColor.clearColor()
        view.addSubview(scrollView)
        scrollView.contentView.addSubview(usernameTextField)
        scrollView.contentView.addSubview(passwordTextField)
        usernameTextField.layer.borderColor = UIColor(rgba: "#4B6C87").CGColor
        usernameTextField.layer.borderWidth
        usernameTextField.placeholder = "账号"
        usernameTextField.attributedPlaceholder = NSAttributedString(string: "账号",
            attributes: [
                NSForegroundColorAttributeName: UIColor(rgba: "#424242"),
            ])
//        usernameTextField.backgroundColor = UIColor.whiteColor()
    }
    
    private func setupConstraints() {
        scrollView.mas_makeConstraints({m in
            m.edges.equalTo()(self.view)
        })
        usernameTextField.mas_makeConstraints({m in
            m.centerX.equalTo()(self.scrollView)
            m.top.mas_equalTo()(400)
        })
        passwordTextField.mas_makeConstraints({m in
            m.centerX.equalTo()(self.scrollView)
        })
    }
}
