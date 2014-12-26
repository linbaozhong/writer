/*
 * accordion 插件 - Ver 0.1
 * by jQuery 2.1.1
 *
 * 哈利蔺特 2014-10-17
 */

(function($) {
	$.fn.accordion = function(options) {
		var $this = this,
			defaults = {
				size: {
					heading: 30, //标题宽度或高度,'auto'可选
					body: 600 //单体宽度或者高度,'auto'可选
				},
				background: 'white',
				horizontal: true, //横向或纵向排列
				clickonly: true, //只允许click触发
				autoopen: 1, //自动打开第几个
				speed: 400 //动画速度
			},
			opts = $.extend(defaults, options),
			frameStyle = {
				'position': 'absolute',
				'left': 5000,
				'top': 0,
				'cursor': 'pointer',
				'overflow': 'hidden',
				'background': opts.background
			},
			getBorder = function(_obj) {
				return {
					'width': parseInt(_obj.css('paddingLeft')) + parseInt(_obj.css('paddingRight')) + parseInt(_obj.css('marginLeft')) + parseInt(_obj.css('marginRight')) + parseInt(_obj.css('borderLeftWidth')) + parseInt(_obj.css('borderRightWidth')),
					'height': parseInt(_obj.css('paddingTop')) + parseInt(_obj.css('paddingBottom')) + parseInt(_obj.css('marginTop')) + parseInt(_obj.css('marginBottom')) + parseInt(_obj.css('borderTopWidth')) + parseInt(_obj.css('borderBottomWidth'))
				};
			},
			refresh = function() {
				//重绘
				$this.each(function() {
					var self = $(this),
						_frame = self.find('.frame'),
						_border = getBorder(_frame);

					// 渲染frame
					if (opts.horizontal) {
						// 横向
						_frame.css({
							'height': self.height() - _border.height,
							'width': opts.size.body - _border.width
						});
					} else {
						// 纵向
						_frame.css({
							'height': (self.height() - ((_frame.length - _frame.filter('.opening').length) * opts.size.heading)) / _frame.filter('.opening').length - _border.height,
							'width': self.width() - _border.width
						});
					}
					// 渲染 heading
					_frame.find('.heading').each(function() {
						var heading = $(this),
							headingBorder = getBorder(heading);
						// 横向
						if (opts.horizontal) {
							heading.css({
								'textAlign': 'center',
								'overflow': 'hidden',
								'float': 'left',
								'width': opts.size.heading - headingBorder.width,
								'height': _frame.height() - headingBorder.height
							});
						} else {
							heading.css({
								'overflow': 'hidden',
								'lineHeight': (opts.size.heading - headingBorder.height) + 'px',
								'height': opts.size.heading - headingBorder.height,
								'width': _frame.width() - headingBorder.width
							});
						}
					});
					// 渲染 content
					_frame.find('.content').each(function() {
						var content = $(this),
							contentBorder = getBorder(content);
						// 横向
						if (opts.horizontal) {
							content.css({
								'width': _frame.width() - contentBorder.width + (content.siblings('.heading').length > 0 ? 0 : opts.size.heading),
								'height': _frame.height() - contentBorder.height,
								'marginLeft': (content.siblings('.heading').length > 0 ? opts.size.heading : 0),
								'overflowX': 'hidden',
								'overflowY': 'auto'
							});
						} else {
							content.css({
								'height': _frame.height() - contentBorder.height - (content.siblings('.heading').length > 0 ? opts.size.heading : 0),
								'width': _frame.width() - contentBorder.width,
								'overflowX': 'hidden',
								'overflowY': 'auto'
							});
						}
					});
				});
			};

		/*
		 * 新抽屉
		 * add({
		 * 	title:'标题',
		 *  html:'html格式编码',
		 *  active:false
		 * })
		 */
		$this.add = function(options) {
			var _options = $.extend({
				title: '',
				html: '',
				url: '',
				active: false
			}, options);

			var _frame = $('<div />').addClass('frame').css(frameStyle).appendTo($this),
				heading = $('<div />').addClass('heading').css('lineHeight', opts.size.heading + 'px').html(_options.title),
				content = $('<div />').addClass('content').html(_options.html);
			//标题
			if (_options.title != '') {
				_frame.append(heading);
			}
			//内容,异步载入
			if (_options.url && _options.url.length > 0) {
				content.load(_options.url);
			}
			_frame.append(content);
			//活动状态
			if (_options.active) {
				_frame.addClass('active').addClass('opening');
			}

			$this.open(_frame);

			refresh();
			return $this;

		};
		//
		$this.active = function(i) {
			var _frame = $this.find('.frame:nth-child(' + i + ')').addClass('active').addClass('opening');
			$this.open(_frame);
			return $this;
		};
		/*
		 * 展开frame
		 */
		$this.open = function(_obj) {
			var _this = _obj.parent(),
				_frame = _this.children(),
				_left = 0,
				_top = 0,
				_width = (_this.width() - (_frame.filter('.opening').length * opts.size.body)) / (_frame.length - _frame.filter('.opening').length),
				_height = (_this.height() - ((_frame.length - _frame.filter('.opening').length) * opts.size.heading)) / _frame.filter('.opening').length;

			_frame.each(function(index) {
				var __this = $(this),
					_border = getBorder(__this);

				if (opts.horizontal) {
					__this.stop().animate({
						'left': _left
					}, opts.speed);

					_left += __this.hasClass('opening') ? __this.width() : _width > _frame.width() ? _frame.width() : _width;
				} else {
					__this.stop().animate({
						'top': _top
					}, opts.speed);
					_top += __this.hasClass('opening') ? _height - _border.height : opts.size.heading;
				}
			});
		};
		//
		$(window).on('resize', function() {
			// 刷新
			$this.each(function() {
				var self = $(this),
					_frame = self.find('.frame.active').addClass('opening');

				self.css({
					height: self.parent().height()
				});

				$this.open(_frame);
				// 刷新
				refresh();
			});
		});
		//
		return $this.each(function() {
			var self = $(this);

			/*
			 * 复原
			 */
			var _revert = function(_obj) {
				if (_obj) {
					self.find('.frame.active').removeClass('active opening');
					_obj.addClass('active opening').next().addClass('active opening');
				} else {
					_obj = self.find('.frame.active').addClass('opening');
				}
				$this.open(_obj);
				$this.open(_obj.next());
			};
			/*
			 * 初始化
			 */
			var _init = (function() {

				//渲染
				self.addClass('accordion').css({
					position: 'relative',
					overflow: 'hidden',
					height: '100%'
				});

				//横向
				if (opts.horizontal) {
					self.addClass('horizontal');
				}

				//自动打开
				var _frame = self.find('.frame').css(frameStyle);
				$this.open(_frame.filter(':nth-child(' + opts.autoopen + ')').addClass('active').addClass('opening'));

				// 刷新
				refresh();

				//事件
				self.on('click', '.frame', function(e) {
					//
					var _obj = $(this);
					//当前活动的frame保持原状
					if (_obj.hasClass('active') && _obj.next().hasClass('active')) {
						return;
					} else if (_obj.next().length == 0) {
						//新建一个抽屉
						$this.add({
							html: 'it.htm',
							active: true
						});
					};

					_revert(_obj);
				}).on('mouseenter', '.frame', function(e) {
					var _that = $(this);
					//当前活动的frame保持原状
					if (_that.hasClass('active')) {
						return;
					}
					self.find('.frame.opening').removeClass('opening');
					_that.addClass('opening');
					$this.open(_that);
				}).on('mouseleave', '.frame', function(e) {
					var _that = $(this);
					//当前活动的frame保持原状
					if (_that.hasClass('active')) {
						return;
					}
					self.find('.frame.opening').removeClass('opening');
					_revert();
				});
			})();

		});

	};

})(jQuery)