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
				space:0,//间距
				during:200,//延时
				timer:null,
				speed: 400 //动画速度
			},
			opts = $.extend(defaults, options),
			frameStyle = {
				'position': 'absolute',
				'overflow': 'hidden',
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
				opts.space = _width - 1;
				// 渲染frame
				_frames.css({
					'height': self.height(),
					'width': opts.size.body
				});
				
				$.each(_frames, function(index, frame) {
					var _frame = $(frame).css({left: _left});

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
					_left += _frame.hasClass('active') ? opts.size.body - 1 : opts.space;
				});
			},

			getWrap = function(){
				return '<div class="accordion-item"></div>'
			},
			_sortable = function(frame) {
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
						snow.log('over');
						// 当前文档属性
						snow.article.parentId = frame.data('parentid')
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
						// 当前文档属性
						snow.article = {
							id : ui.item.data('id'),
							parentId : ui.item.data('parentid'),
							documentId : ui.item.data('documentid')
						}

						ui.item.show();
						// 如果当前节点的子级已经是编辑状态,禁止拖拽
						var _child = ui.item.closest('.frame').next('.frame');
						
						if(_child.length && ui.item.data('id')===_child.data('parentid')){
							frame.sortable('cancel');
						}else{
							// 如果是作者的作品,可以任意拖拽,否则，只能克隆
							if(snow.author(ui.item)){
								ui.item.hide();
							}else{
							}
						}	
					},
					beforeStop:function(e,ui){
						// 如果是作者的作品,可以任意拖拽,否则，只能克隆
						var _doc;
						if(snow.author(ui.item)){
							_doc = ui.item;
						}else{
							_doc = ui.item.clone();
							ui.placeholder.after(_doc.data('id','0'));
							frame.sortable('cancel');
						}
					},
					stop:function(e,ui){
						
					},
					out: function(e, ui) {
						//clearTimeout(opts.timer);	
					},
					change:function(e,ui){
						snow.log('change');
						// 当前文档属性
						snow.article.position = ui.placeholder.prev().data('id');
					}
				}).droppable({
					accept:'div.doc',
					drop:function(e,ui){
						//snow.log(ui.draggable.data('parentid'),ui.draggable.prev('div.doc').data('id'),snow.article);
						// 如果位置发生变化
						if(ui.draggable.data('parentid') != snow.article.parentId 
							|| ui.draggable.prev('div.doc').data('id') != snow.article.position){
							// 如果是作者的作品,可以任意拖拽,否则，只能克隆
							if(snow.author(ui.draggable)){
								// 只修改parentId和position
								$.post(snow.api.docPosition,{
									id : snow.article.id,
									parentId : snow.article.parentId,
									documentId : snow.article.documentId,
									referId : snow.article.position
								},function(result){
									snow.log(result);
								});					
									
							}else{
								// 新建一个article，documentId不变
							};
						};
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
				parentId:''
			}, options);

			var _frame = $('<div />').attr('id',Math.random()).addClass('frame').css(frameStyle).css({
				left: 5000
			}).appendTo($this).data('parentid',_options.parentId);
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
			
			if (parseInt(_obj.css('left')) > _obj.index()* opts.space) {
				//console.log('after');
				_frames = _obj.prevAll();
				$.each(_frames,function(){
					var __frame=$(this);

					__frame.stop().animate({
							left: __frame.index() * opts.space
					}, opts.speed);
				});
				_obj.stop().animate({'left':_obj.index() * opts.space}, opts.speed);
			}else{
				//console.log('before');
				_frames = _obj.nextAll();
				$.each(_frames,function(){
					var __frame=$(this);

					__frame.stop().animate({
							left: (__frame.index()-1) * opts.space + opts.size.body-1
					}, opts.speed);
				});
			}
			_obj.addClass('active').siblings('.active').removeClass('active');
		};
		//
		$(window).on('resize', function() {
			// 刷新
			$this.each(function() {
				var self = $(this);

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
				var _frames = self.find('.frame').css(frameStyle);

				$.each(_frames, function() {
					var __frame = $(this);
					// 设置id属性
					if (!__frame.attr('id')) {
						__frame.attr('id',Math.random());
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
					opts.timer = setTimeout(function(){
							$this.open(_that);
						},opts.during);
				}).on('mouseleave','.frame',function(e){
					clearTimeout(opts.timer);					
				});
			})();

		});

	};

})(jQuery);
