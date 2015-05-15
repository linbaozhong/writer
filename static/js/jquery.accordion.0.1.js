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
				space: 0, //间距
				during: 200, //延时
				timer: null,
				speed: 400 //动画速度
			},
			opts = $.extend(defaults, options),
			frameStyle = {
				'position': 'absolute',
				'background': opts.background
			},
			//重绘frame
			refresh = function(self) {
				var _left = 0,
					_frames = self.find('.frame'),
					_width = (self.width() - opts.size.body) / (_frames.length - 1);
				// 
				_width = _width > opts.size.body ? opts.size.body : _width;
				// 间距
				opts.space = _width;
				// 渲染frame
				_frames.css({
					'minHeight': self.height(),
					'width': opts.size.body
				});

				$.each(_frames, function(index, frame) {
					var _frame = $(frame),
						_height = self.height(),
						_maxHeight = _frame.outerHeight(),
						_top = _frame.position().top;
						
					if (_maxHeight + _top < _height) {
						_top = _height - _maxHeight;
					}
					_frame.css({
						left: _left,
						top:_top
					});

					// 遮罩
					_frame.find('.mask').css({
						width: _frame.outerWidth(),
						height: _frame.outerHeight()
					});

					if (_frame.position().left != _left) {
						_frame.stop().animate({
							left: _left
						}, opts.speed);
					}
					
					if(_frame.hasClass('active')){
						// 刷新滚动条
						snow.autoScroll(_frame);
						
						_left += opts.size.body;
					} else{
						_left += opts.space;
					}
				});
			},

			getWrap = function() {
				return '<div class="accordion-item"></div>'
			},
			_sortable = function(frame) {
//				frame.droppable({
//					drop:function(e,ui){
//						console.log(ui.draggable,ui.helper);
//					}
//				});
				frame.sortable({
					connectWith: 'div.frame',
					appendTo: '#snow-left',
					helper: 'clone',
					items: 'div.doc',
					handle: ".doc-handle",
					placeholder: "portlet-placeholder ui-corner-all",
					//revert: true,
					tolerance: 'pointer',
					opacity: 0.8,
					over: function(e, ui) {
						snow.article.frame = frame;
						// 当前文档属性
						snow.article.moreId = frame.data('moreid')
						snow.article.parentId = frame.data('parentid')

						// 能否拖拽，当前激活文档不能向右拖，并且，已经存在的文档不能拖入
						snow.article.disable = (ui.item.hasClass('active') && (frame.index() > ui.item.closest('.frame').index())) 
							|| ((frame.index() != ui.item.closest('.frame').index()) && frame.find('#' + ui.item.attr('id')).length)
							|| (snow.article.parentId == ui.item.data('id'));
						//当前活动的frame保持原状
						if (frame.hasClass('active')) {
							return;
						}

						$this.open(frame);

						if (ui.placeholder.prevAll('footer').length) {
							ui.placeholder.after(ui.placeholder.prevAll('footer'));
						}
					},
					start: function(e, ui) {
						// 记录文档的参考位置id(前一个文档的id)
						ui.item.data('referid', ui.item.prev('div.doc').data('moreid'));
						// 当前文档属性
						snow.article = {
							id: ui.item.data('moreid'),
							parentId:0,
							disable: true
						};

						// 如果是作者的作品,可以任意拖拽,否则，只能克隆
						if (snow.updator(ui.item)) {
							ui.item.hide();
						} else {
							ui.item.show();
						};

					},
					beforeStop: function(e, ui) {
						// 防止自为父节点拖拽和重复文档，并且只能对自己的frame操作
						if ((!snow.article.frame.hasClass('snow-me') && !snow.updator(ui.item)) || snow.article.disable) {
							frame.sortable('cancel');
							return;
						}
						// 如果是作者的作品,可以任意拖拽,否则，只能克隆
						var _doc;
						if (snow.updator(ui.item)) {
							_doc = ui.item;
						} else {
							_doc = ui.item.clone();
							ui.placeholder.after(_doc.data('id', '0'));
							frame.sortable('cancel');
						}

						if (_doc.data('parentid') != snow.article.parentId) {
							// 清除激活状态
							_doc.removeClass('active');
						}

						// 如果位置发生变化
						if (_doc.data('parentid') != snow.article.parentId || _doc.data('referid') != snow.article.referId) {
							// 记录新位置
							_doc.data('parentid', snow.article.parentId).data('referid', snow.article.referId)
							// 
							snow.article.doc = _doc;
							// 只修改parentId和position
							$.post(snow.api.docPosition, {
								id: snow.article.id, //文档的moreId
								moreId: snow.article.moreId, //新的父节点
								referId: snow.article.referId //参考节点
							}, function(result) {
								snow.log(result);
								if(result.ok){
									var _more = result.data;
									snow.article.doc.data('moreid',_more.id).data('updator',_more.updator);
								}
							});
							
						};
					},
					sort: function(e, ui) {
						// 当前文档属性
						snow.article.referId = ui.placeholder.prev('div.doc').data('moreid');
					}
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
		$this.add = function(options, fn) {
			var _options = $.extend({
				title: '',
				html: '',
				url: '',
				active: false,
				moreId: ''
			}, options);

			var _frame = $('<div />').attr('id', Math.random()).addClass('frame').css(frameStyle).css({
				left: 5000
			}).appendTo($this).data('moreid', _options.moreId);
			//活动状态
			if (_options.active) {
				_frame.addClass('active');
			}

			//
			if ($.isFunction(fn)) {
				fn(_frame);
			}
			//
			_sortable(_frame);

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
			var _frames;

			if (parseInt(_obj.css('left')) > _obj.index() * opts.space) {
				//console.log('after');
				_frames = _obj.prevAll();
				$.each(_frames, function() {
					var __frame = $(this);

					__frame.stop().animate({
						left: __frame.index() * opts.space
					}, opts.speed);
				});
				
				_obj.stop().animate({
					'left': _obj.index() * opts.space
				}, opts.speed,function(){
					// 刷新滚动条
					snow.autoScroll(_obj);
				});
			} else {
				// 刷新滚动条
				snow.autoScroll(_obj);
				
				_frames = _obj.nextAll();
				$.each(_frames, function() {
					var __frame = $(this);

					__frame.stop().animate({
						left: (__frame.index() - 1) * opts.space + opts.size.body
					}, opts.speed);
				});
			}
			_obj.addClass('active').siblings('.active').removeClass('active');
			
		};
		// 刷新重绘
		$this.refresh = function() {
			refresh($(this));
		};
		//
		$(window).on('resize', function() {
			// 刷新
			$this.each(function() {
				var self = $(this);
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
				var _frames = self.find('.frame').css(frameStyle);

				$.each(_frames, function() {
					var __frame = $(this);
					// 设置id属性
					if (!__frame.attr('id')) {
						__frame.attr('id', Math.random());
					}
					_sortable(__frame);
				});
				//重置
				refresh(_frames.eq(0).parent());

				$this.open(_frames.filter(':nth-child(' + opts.autoopen + ')'));


				//事件
				self.on('mouseenter', '.frame', function(e) {
					var _that = $(this);
					//当前活动的frame保持原状
					if (_that.hasClass('active')) {
						return;
					}
					opts.timer = setTimeout(function() {
						$this.open(_that);
					}, opts.during);
				}).on('mouseleave', '.frame', function(e) {
					clearTimeout(opts.timer);
				});
			})();

		});

	};

})(jQuery);