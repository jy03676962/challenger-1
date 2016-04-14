//
//  CustomNotification.swift
//  postgame
//
//  Created by tassar on 4/14/16.
//  Copyright © 2016 pulupulu. All rights reserved.
//

import Foundation

public enum NotificationKey: String {
	case HostChanged = "HostChangedNotification"
}

extension NSNotificationCenter {

	func addObserver(observer: AnyObject, selector aSelector: Selector, key aKey: NotificationKey) {
		self.addObserver(observer, selector: aSelector, name: aKey.rawValue, object: nil)
	}

	func addObserver(observer: AnyObject, selector aSelector: Selector, key aKey: NotificationKey, object anObject: AnyObject?) {
		self.addObserver(observer, selector: aSelector, name: aKey.rawValue, object: anObject)
	}

	func removeObserver(observer: AnyObject, key aKey: NotificationKey, object anObject: AnyObject?) {
		self.removeObserver(observer, name: aKey.rawValue, object: anObject)
	}

	func postNotificationKey(key: NotificationKey, object anObject: AnyObject?) {
		self.postNotificationName(key.rawValue, object: anObject)
	}

	func postNotificationKey(key: NotificationKey, object anObject: AnyObject?, userInfo aUserInfo: [NSObject: AnyObject]?) {
		self.postNotificationName(key.rawValue, object: anObject, userInfo: aUserInfo)
	}

	func addObserverForKey(key: NotificationKey, object obj: AnyObject?, queue: NSOperationQueue?, usingBlock block: (NSNotification!) -> Void) -> NSObjectProtocol {
		return self.addObserverForName(key.rawValue, object: obj, queue: queue, usingBlock: block)
	}
}