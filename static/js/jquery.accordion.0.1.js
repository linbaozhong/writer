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
				clickonly: true, //只允许click触发
				autoopen: 1, //自动打开第几个
				speed: 400 //动画速度
			},
			opts = $.extend(defaults, options),
			frameStyle = {
				'position': 'absolute',
				'cursor': 'pointer',
				'overflow': 'hidden',
				'background': opts.background
			},
			getBorder = function(_obj) {
				return {
					'width': parseInt(_obj.css('paddingLeft')) + parseInt(_obj.css('paddingRight')) + parseInt(_obj.css('borderLeftWidth')) + parseInt(_obj.css('borderRightWidth')),
					'height': parseInt(_obj.css('paddingTop')) + parseInt(_obj.css('paddingBottom')) + parseInt(_obj.css('borderTopWidth')) + parseInt(_obj.css('borderBottomWidth'))
				};
			},
			//重绘
			refresh = function(self) {
				var _left = 0,
					_frames = self.find('.frame'),
					_width = (self.width() - opts.size.body) / _frames.length - 1;

				// 
				_width = _width > opts.size.body ? opts.size.body : _width;

				// 渲染frame
				_frames.css({
					'height': self.height(),
					'width': opts.size.body
				});

				$.each(_frames, function(index, frame) {
					var _frame = $(frame);

					if (_frame.position().left != _left) {
						_frame.stop().animate({
							left: _left
						}, opts.speed);
					}
					_left += _frame.hasClass('active') ? opts.size.body - 1 : _width - 1;
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

			var _frame = $('<div />').addClass('frame').css(frameStyle).css({
				left: 5000
			}).appendTo($this);
			//活动状态
			if (_options.active) {
				_frame.addClass('active');
			}
			_frame.sortable({
				connectWith: 'div.frame',
				appendTo: '#snow-left',
				items:'div.doc',
				//helper: 'clone',
				handle: ".doc-info",
				placeholder: "portlet-placeholder ui-corner-all",
				revert: true,
				tolerance:'pointer',
				opacity: 0.8,
				over:function(e,ui){
					_frame.closest('#snow-left').find('.frame.active').removeClass('active');
					$this.open($(e.target));
				},
				start: function(e, ui) {
				},
				stop: function(e, ui) {
				}
			});
			refresh(_frame.parent());
			return $this;

		};
		//
		$this.active = function(i) {
			var _frame = $this.find('.frame:nth-child(' + i + ')').addClass('active');
			$this.open(_frame);
			return $this;
		};
		/*
		 * 展开frame
		 */
		$this.open = function(_obj) {
			_obj.addClass('active');

			refresh(_obj.parent());
		};
		//
		$(window).on('resize', function() {
			// 刷新
			$this.each(function() {
				var self = $(this),
					_frame = self.find('.frame.active');

				self.css({
					height: self.parent().height()
				});

				// 刷新
				refresh(self);
			});
		});
		//
		return $this.each(function() {
			var self = $(this);

			/*
			 * 复原
			 */
			var _revert = function() {
				// 刷新
				refresh(self);
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
				self.addClass('horizontal');

				//自动打开
				var _frame = self.find('.frame').css(frameStyle);
				$this.open(_frame.filter(':nth-child(' + opts.autoopen + ')').addClass('active'));


				//事件
				self.on('click', '.frame div.doc', function(e) {
					//
					self.find('.frame.active').removeClass('active');
					var _obj = $(this).closest('.frame');

					if (_obj.next().length == 0) {
						//新建一个抽屉
						$this.add({
							html: 'it.htm',
							active: true
						});
					} else {
						_obj.next().addClass('active');
					};
					_revert();
				}).on('mouseenter', '.frame', function(e) {
					var _that = $(this);
					//当前活动的frame保持原状
					if (_that.hasClass('active')) {
						return;
					}
					self.find('.frame.active').removeClass('active');
					$this.open(_that);
				}).on('mouseleave', '.frame', function(e) {
					var _that = $(this);
					//当前活动的frame保持原状
					if (_that.hasClass('active')) {
						return;
					}
					_revert();
				});
			})();

		});

	};

})(jQuery)