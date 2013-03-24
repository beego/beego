//scripts for wac

var tvInit = function(){
     var currentTallest = 0,
     currentRowStart = 0,
     rowDivs = new Array(),$el,topPosition = 0;
  
   $('.videoList li').each(function() {
  
     $el = $(this);
     topPostion = $el.position().top;
  
     if (currentRowStart != topPostion) {
  
     // we just came to a new row.  Set all the heights on the completed row
     for (currentDiv = 0 ; currentDiv < rowDivs.length ; currentDiv++) {
       rowDivs[currentDiv].height(currentTallest);
     }
  
     // set the variables for the new row
     rowDivs.length = 0; // empty the array
     currentRowStart = topPostion;
     currentTallest = $el.height();
     rowDivs.push($el);
  
     } else {
  
     // another div on the current row.  Add it to the list and check if it's taller
     rowDivs.push($el);
     currentTallest = (currentTallest < $el.height()) ? ($el.height()) : (currentTallest);
  
    }
  
    // do the last row
     for (currentDiv = 0 ; currentDiv < rowDivs.length ; currentDiv++) {
     rowDivs[currentDiv].height(currentTallest);
     }
  
   });
  
  $('#searchBox #searchButton,#gridView').click(function() {
    if($('#videoPlayer').is(':visible') == true){
      hideVideoPlayer();
    }
  });
  
  $('.videoList,#videoView').click(function() {
    if($('#videoPlayer').is(':hidden') == true){
      showVideoPlayer();
    }
  });
  function showVideoPlayer(){
    $('#videoPlayer').slideDown('fast');
    $("#videoView").addClass("selected");
    $("#gridView").removeClass("selected");
  }
  function hideVideoPlayer(){
    $('#videoPlayer').slideUp('fast');
    $("#gridView").addClass("selected");
    $("#videoView").removeClass("selected");
  }    
}

var productTours = function(){     
  var History = window.History; // Note: We are using a capital H instead of a lower h
  //first make arrays of all urls to fetch and div names from the rel tag
  var contentRels = [];

  $('#featureNav a.tourLink').each(function(i){
      contentRels.push($(this).attr('rel'));
  });

    //create empty divs with the correct div (from the rel names on the links) id's for populating with ajax content 
  for (var i=0; i < contentRels.length; i++) {                         
    if ($('.featureSection').attr('id') !== $(contentRels)[i]) {
        
       $('#tourContent .inner').append('<div id=' + $(contentRels)[i] + ' class="featureSection"><div style="height:1000px; padding: 30px; background-position: center 50px;" class="loadingAjax"></div></div>'); 
    };
    };

  var arrayOfUrls = [];
  var arrayOfDivs = [];  
  

  $('.featureSection').each(function(){
    var divId = $(this).attr('id');
    var urlToGetContentFrom =  $('#featureNav a[rel='+'"'+divId+'"]').attr('href');      
    arrayOfDivs.push(divId);
    arrayOfUrls.push(urlToGetContentFrom);

  });

  
  function loadAjax(){ 

    if (arrayOfUrls.length > 0) {

      var divToPopulate = arrayOfDivs.shift();
    var urlOfContent =  arrayOfUrls.shift();

      if (!$('#'+divToPopulate).hasClass('onLoadContent')) {
        
        $('#'+divToPopulate).load(urlOfContent + ' .featureSectionInner', function(){});

      }
      loadAjax();
    }

  }                      

  loadAjax();
         
    // $('.featureSection').not($('#screenshot-tour')).hide(); 

    // $('#screenshot-tour').css({'visibility' : 'hidden', 'position':'absolute', 'display': 'block'});
        

  var divToShow = $('a[href="'+window.location.pathname+'"]').attr('rel');
    showContent(divToShow);
    if (divToShow) {              
      $(window).scrollTo('#featureTour');
    }
      
     
    // Bind to State Change
    History.Adapter.bind(window,'statechange',function(){ // Note: We are using statechange instead of popstate

        var State = History.getState(); // Note: We are using History.getState() instead of event.state
        var divToShow = State.data.divToShow;

        $('#featureNav a').closest('li').removeClass('selected');
        $('a[rel="'+divToShow+'"]').closest('li').addClass('selected');
        showContent(divToShow);              

        // $('title').html(history.state.titleNode);
    });



  if ($.browser.msie) {
    window.statechange = function(){
      $(window).scrollTo($("#featureTour").offset().top, 0);     
      }
  }

             
   $('#featureNav a').click(function(e){
     
      if ($(this).hasClass('tourLink')) {
        e.preventDefault();
                 
        var contentNode = $(this).attr('rel');
        var fakeLink = $(this).attr('href');
        var titleTag = $("#"+contentNode).find('.featureSectionInner').attr('rel') + " | Atlassian";             
        History.pushState({divToShow : contentNode}, titleTag, fakeLink);
            
        $('#featureNav a').closest('li').removeClass('selected');
        $('a[rel="'+contentNode+'"]').closest('li').addClass('selected');
           
      }
      
      
      
  if ($('.stickyNav').hasClass('stuck')) {                
    var headingSize = $("#featureTour h2").outerHeight(true);
    if ($(this).hasClass('tourLink')) {
      $(window).scrollTo($("#featureTour").offset().top + headingSize, 400, {axis : 'y'});   
    }
  }
    

   });

      function showContent(divToShow, transition){ 
       
          
    $('a[rel="'+divToShow+'"]').closest('li').addClass('selected');
    var divToShowEl = $('#' + divToShow);     


       // console.log(divToShowEl.length);  


      if (divToShowEl.length === 0) {        
    
    if (!$('#featureNav li').hasClass('selected')) {
      $('#featureNav a:first').closest('li').addClass('selected');            
    }
    
    divToShowEl = $('.featureSection:first'); 
            
      } else {
            $('a[rel="'+divToShow+'"]').closest('li').addClass('selected');
    }
        


      
    footerPageNav();
    
    if ($('#featureNav li:first').hasClass('selected')) {
    $('#tourContent').addClass('navTop');

    } else {
    $('#tourContent').removeClass('navTop');
    }    
    

    if (transition) {
      var fading = false;

      if (fading) {
        $('.featureSection:visible').fadeOut(200, function(){
          divToShowEl.css({'visibility' : 'visible', 'position':'static'});
        // $('.featureSection:visible').hide();
          divToShowEl.fadeIn(200);     

          var divHeight = $('.featureSection:visible').outerHeight();

          // if (divToShow.attr('id') == "screenshot-tour") {
          //   divHeight = divHeight +30;
          // }

          // $('#tourContent').animate({'height':divHeight}, 500);

        });
      } else {
        $('.featureSection:visible').hide();
          divToShowEl.css({'visibility' : 'visible', 'position':'static'});


          divToShowEl.show();     

          var divHeight = $('.featureSection:visible').outerHeight();

          // if (divToShow.attr('id') == "screenshot-tour") {
          //   divHeight = divHeight +30;
          // }            

          // $('#tourContent').animate({'height':divHeight}, 1000, function(){
          // });                      

      }

    } else {
   
      $('.featureSection:visible').hide();                  
        divToShowEl.css({'visibility' : 'visible', 'position':'static'})
        divToShowEl.show();

    }

  }     
     
      
  $('html').live("keydown", function(e) {                                          
        var headingSize = $("#featureTour h2").outerHeight(true);
    var nextContentLink = $('#featureNav .selected').next('li').find('.tourLink').attr('rel');
    var prevContentLink = $('#featureNav .selected').prev('li').find('.tourLink').attr('rel');
    var prevSSItem = $('ul.ssList').find('li.selected').prev();
    var nextSSItem = $('ul.ssList').find('li.selected').next();
    var fakeLink = $('#featureNav .selected').next('li').find('.tourLink').attr('href');
    var prevFakeLink = $('#featureNav .selected').prev('li').find('.tourLink').attr('href');    


    var nextTitleTag = $("#"+nextContentLink).find('.featureSectionInner').attr('rel');
    var prevTitleTag = $("#"+prevContentLink).find('.featureSectionInner').attr('rel');
    // $('.ssNext').click();        
    //         $('.ssPrev').click();
    if (e.which == 39) {                   
        if (nextContentLink) {
          History.pushState({divToShow : nextContentLink}, nextTitleTag, fakeLink);          
          
          $(window).scrollTop(($("#featureTour").offset().top + headingSize));
          
        } else if ($('#full-screenshot-tour').is(':visible') && nextSSItem.length > 0){
          $('.ssNext').click();
        } 
      }
      

      if (e.which == 37) {
        
        if ($('#full-screenshot-tour').is(':visible') && prevSSItem.length > 0){
          $('.ssPrev').click();
        } else {
          if (prevContentLink) {
            History.pushState({divToShow : prevContentLink}, prevTitleTag, fakeLink);          
            $(window).scrollTop(($("#featureTour").offset().top + headingSize));
          }
        }

        }

                   
  });
  
  $('.next').live('click', function(e) {
        var headingSize = $("#featureTour h2").outerHeight(true);
    var nextContentLink = $('#featureNav .selected').next('li').find('.tourLink').attr('rel');
    var prevContentLink = $('#featureNav .selected').prev('li').find('.tourLink').attr('rel');
    var fakeLink = $('#featureNav .selected').next('li').find('.tourLink').attr('href');
    var prevFakeLink = $('#featureNav .selected').prev('li').find('.tourLink').attr('href');    
    var nextTitleTag = $('#featureNav .selected').next('li').find('.tourLink').html();
  nextTitleTag = nextTitleTag.replace('&amp;', 'and') + " | Atlassian";
 
    e.preventDefault();
    if (nextContentLink) {

      History.pushState({divToShow : nextContentLink}, nextTitleTag, fakeLink);                
      $(window).scrollTop(($("#featureTour").offset().top + headingSize));
    }
  });           
  
  $('.previous').live('click', function(e) {
        var headingSize = $("#featureTour h2").outerHeight(true);
    var nextContentLink = $('#featureNav .selected').next('li').find('.tourLink').attr('rel');
    var prevContentLink = $('#featureNav .selected').prev('li').find('.tourLink').attr('rel');
    var fakeLink = $('#featureNav .selected').next('li').find('.tourLink').attr('href');
    var prevFakeLink = $('#featureNav .selected').prev('li').find('.tourLink').attr('href');    
    var previousTitleTag = $('#featureNav .selected').prev('li').find('.tourLink').html(); 
  previousTitleTag = previousTitleTag.replace('&amp;', 'and') + " | Atlassian";
    e.preventDefault();
    if (prevContentLink) {
      History.pushState({divToShow : prevContentLink}, previousTitleTag, prevFakeLink);
      $(window).scrollTop(($("#featureTour").offset().top + headingSize));
    }    
  });
  

};

var featureStickyProductNavInit = function(){

    var scrollBottom = $(window).scrollTop() + $(window).height();
    var documentHeight = $(document).height();
    var footerHeight = $('#footer').height(); //this will need to include the product footer as well

  $(window).bind("resize", $.throttle( 10, function(){
     $(window).scroll();
      // if ($('#featureNav').height() < $(window).height() && $(window).scrollLeft() == 0) {
      // stickNavigation(); 
      // }   
  
  }));
         

         if ($('#featureNav').height() < $(window).height() && $(window).scrollLeft() == 0) {
        stickNavigation(); 
  }

  function stickNavigation() {
    
    //footer lock position (footer height) is the distance of the footer from the top minus the documnet height 
    
    
      var headerHeight = $('#featureNav').offset().top;
      var navPosition = $('#featureNav').offset().top;
      var navOffsetPosition = $('.stickyNav').offset().top;
      var footerPosition = $("#endOfTour").offset().top;
      var lastItemPosition = $('.navClone').offset().top + ($('.navClone').outerHeight() + 30);

      $('#featureList').css('position', 'relative');

      $(".navClone").css({
          'height': $('#featureNav').height(),
          'width': $('#featureNav').outerWidth(),
          'visibility' : 'hidden',
          'margin-left' : -$('#featureNav').outerWidth(),
          'float' : 'left'
      });

   

    $(window).bind("scroll", $.throttle( 0, function(){
      navPosition = $('#featureNav').offset().top;
          navOffsetPosition = $('.stickyNav').offset().top;
          footerPosition = $("#endOfTour").offset().top - 50;
          lastItemPosition = parseInt($('.navClone').offset().top + ($('.navClone').height()));
          originalItemPosition = parseInt($('#featureNav').offset().top + ($('#featureNav').height()));
      
          setNavPosition({
              headerHeight : headerHeight,
              navPosition : navPosition,
              navOffsetPosition : navOffsetPosition,
              footerPosition : footerPosition,
              lastItemPosition : lastItemPosition
          });
    
    if (($(window).width() < 960)) {
      if ($('.stickyNav').hasClass('stuck')) {
        $('#featureNav li.selected').css({'zIndex' : '0'});
      }
    } else {
      $('#featureNav li.selected').css({'z-index' : '10'});
    }


          
    }));

      // $(window).scroll(function() {
      // 
      // });  

  }


  function setNavPosition(elValues) {
      var scrollBottom = $(window).scrollTop() + $(window).height();
      var windowHeight = $(window).height();
      var documentHeight = $(document).height();


      if (elValues.lastItemPosition < elValues.footerPosition) {
          if ($(document).scrollTop() > (elValues.headerHeight)) {
      $('.stickyNav').css({'position' : 'fixed', 'top' : '0', 'z-index' : '1000'}).addClass('stuck');
          } else {
      $('.stickyNav').css({'position' : 'static'}).removeClass('stuck')
          }
      } else {
          $('#featureNav').css({'position' : 'absolute', 'bottom' : '68px', 'top' : 'auto'});
      }

  }



}; 
   
var footerPageNav = function(){ 
    var headingSize = $("#featureTour h2").outerHeight(true);
    var nextContentTitle = $('#featureNav .selected').next('li').find('.tourLink').text();
    var prevContentTitle = $('#featureNav .selected').prev('li').find('.tourLink').text();

  var titleTag = $('#featureNav .selected').next('li').find('.tourLink').html(); //THIS IS FOR NOW!
    function checkPrevNext(){ 
         if (prevContentTitle) { 

        $('.prevText').text(prevContentTitle);
        $('.previous').show();

      } else {  
        
        $('.previous').hide();
      }

      if (nextContentTitle) {
        $('.nextText').text(nextContentTitle);
        $('.next').show();
      } else {
        $('.next').hide();
      }
  }
  

                          
  checkPrevNext();
                              
};

var screenshotTour = function(){


  var animationSpeed = 200;

  $('.screenshotTourWrap').fadeIn('fast');
  // console.log($('.screenshotTourWrap'));
  $('.featureScreenshot .ssTN,.featureScreenshot .ssOuter').fadeTo(100,1);  
  $('.ssList li:first').addClass('selected');

  $('.ssList a').click(function(e){
    e.preventDefault();
    //reset all from slider nav thing
    // $('.featureScreenshot').css({'left': '0', 'right': '0'});
    var slideToFadeIn = $(this).attr('id').split('-')[1] - 1;
    $('.ssList li').removeClass('selected');
    $(this).closest('li').addClass('selected');
    
    if (!$('.featureScreenshot').eq(slideToFadeIn).hasClass('active')) {
      $('.featureScreenshot:visible').fadeOut('500').removeClass('active');             
      $('.featureScreenshot').eq(slideToFadeIn).fadeIn('500').addClass('active');
      checkNextPrev();
    }
      var headingToShow = slideToFadeIn + 1;
    showHeading(headingToShow);
        //maybe remove the keyboard navigation on thing from here
  });
  
  $('.slideHeadingItem:first').show();
  
  $('.featureScreenshot:first').addClass('active');
  
  function showHeading(num){
     $('.slideHeadingItem').hide();
     $('#slideHeading-' + num).show();
  }

checkNextPrev();
  

  $('.ssNext').click(function(e){
        e.preventDefault();
    var nextItem = $('ul.ssList').find('li.selected').next();
    var headingToShow = nextItem.attr('id').split('-')[1];
    var divToShow = headingToShow -1;
    if (nextItem.length > 0) {
      $('.ssList a').parent().removeClass('selected');      
      nextItem.addClass('selected');
    }

    showHeading(headingToShow);
    slideLeft(divToShow);
    checkNextPrev();
  }); 
  

  $('.ssPrev').click(function(e){
        e.preventDefault();
    var prevItem = $('ul.ssList').find('li.selected').prev();
    var headingToShow = prevItem.attr('id').split('-')[1];
    var divToShow = headingToShow;

    if (prevItem.length > 0) {
      $('.ssList a').parent().removeClass('selected');      
      prevItem.addClass('selected');
    }                  

    showHeading(headingToShow);
    slideRight(divToShow);
    checkNextPrev();
  });
 
  function checkNextPrev(){
    var currentSlide = $('ul.ssList').find('li.selected');
    var nextSlide = currentSlide.next();
    var previousSlide = currentSlide.prev();
                
    if (nextSlide.length == 0) {
       $('.ssNext').hide(); 
    } else {
      $('.ssNext').show(); 
    }

    if (previousSlide.length == 0) {
       $('.ssPrev').hide(); 
    } else {
      $('.ssPrev').show(); 
    }


  }  
  
  function slideLeft(num){        
    var previousSlideNum = num -1;
    $('.featureScreenshot').removeClass('active');
    
    $('.featureScreenshot').eq(num).css({'right': '-617px', 'display' : 'block'}).addClass('active').animate({'right': '0'}, animationSpeed, function(){
      $(this).css({'right': '0'});
    }); 
    $('.featureScreenshot').eq(previousSlideNum).animate({'right': '617px'}, animationSpeed, function(){
      $(this).css({'display': 'none', 'right' : 'auto'});
    });
     checkNextPrev(); 

  }
  
  function slideRight(num){
        var previousSlideNum = num -1;
    $('.featureScreenshot').removeClass('active');
    
    $('.featureScreenshot').eq(previousSlideNum).css({'right': '617px', 'display' : 'block'}).addClass('active').animate({'right': '0'}, 300, function(){
      $(this).css({'right': '0'});
    });    
    $('.featureScreenshot').eq(num).animate({'right': '-617px'}, 300, function(){
      $(this).css({'display': 'none', 'right' : 'auto'});
    });
     checkNextPrev();                            
     
  }   
  
  var windowHeight = $(window).height();
  var windowWidth = $(window).width();    
  $(window).resize(function(){
    windowHeight = $(window).height();
    windowWidth = $(window).width();
    $("a[rel='[slideShowLB]']").colorbox({
        maxHeight: 1000
        // maxWidth: windowWidth  
    });
  });     

  $("a[rel='[slideShowLB]']").colorbox({
      maxHeight: 1000,
      // maxWidth: windowWidth,  
      current:false,
      transition: 'elastic'

  }); 
  
  //      var highestItem = 70;
  // $('.slideHeadingItem').each(function(i){
  //   if ($('.slideHeadingItem').eq(i).height() > highestItem) {
  //     highestItem = $('.slideHeadingItem').eq(i).height();
  // 
  //   }
  // });                                         
  // console.log(highestItem);
  // setTimeout("$('.slideHeadingItem').css('height', highestItem)", 200);  
};

     

       
var simpleTabNavigation  = function(options){

  var settings = $.extend( {
    listItem: ".nav-section",
    contentItem: ".section",
    bookmarkable: true, 
    scroll: false,
    paramName: 'tab'
      }, options);

      var  queryString =  $.url().param(settings.paramName);
        
    // console.log($.url().param());

    if (queryString) {                                                            
      chooseTabfromParam(queryString);
    } else {
          $(settings.contentItem).not(':first').hide();
          $(settings.listItem + ':first').addClass('selected');
      }
        if (settings.bookmarkable) {
          History.Adapter.bind(window,'statechange',function(){ // Note: We are using statechange instead of popstate
          var State = History.getState(); // Note: We are using History.getState() instead of event.state
          var historyParam = State.data.tabParam;
          if (undefined != historyParam) {
            chooseTabfromParam(historyParam);      
          }
           });
        }
     
  
  function chooseTabfromParam(tabParamName){
    var navItemLinks = $(settings.listItem).find('a');
    var activeItem = $(settings.listItem).find('a[href="#'+tabParamName+'"]'); 
    var positionOfHashItemDiv = navItemLinks.index(activeItem);                
        if (positionOfHashItemDiv !== -1) {
            $(settings.contentItem).not(activeItem).hide();
        $(settings.listItem).removeClass('selected');                       
        $(settings.listItem).eq(positionOfHashItemDiv).addClass('selected');
        showContent(positionOfHashItemDiv);
        }

  }


  if (settings.scroll && queryString) {
    $.scrollTo($('#'+ queryString), 1000);
  }
    


  
  $(settings.listItem).each(function(i){
    $(this).click(function(e){
      e.preventDefault();
      $(settings.listItem).removeClass('selected');
      $(this).addClass('selected');

      if ($(this).find('a').length > 0) {
    var tabName = $(this).find('a').attr('href').split('#')[1]; 
    } else {
      var tabName = $(this).attr('id');     
  }
      
    showContent(i); 
      removeCurveFirstItem(); 
  
      if (settings.bookmarkable == true) {
    //This tabParam state needs to inclue all of the params not just tab Name
        History.replaceState({tabParam: tabName}, $('title').html(), "?" + settings.paramName + "="+ tabName);          
  }
    });        
   });


  function showContent(i){ 
    $(settings.contentItem).hide();
    $(settings.contentItem).eq(i).show();
  // $(settings.contentItem).eq(i);
  }
  
  removeCurveFirstItem();
};

var searchableTabNavigation  = function(settings){

    if (!settings) {
       var settings = {
         listItem: ".nav-section",
         contentItem: ".section",
         scroll: false
       }
    }     
    
    var  queryString =  location.hash;
    
  if (queryString) {                                                            
        showContent("" + queryString);
  } else {
        showContent($(settings.listItem + ':first').attr("name"));
  }
  
    $(settings.listItem).each(function(i){
        $(this).click(function(e){
          showContent($(this).attr("name"));
        });        
    });
    
    var currentState = "";

    History.Adapter.bind(window,'hashchange',function(){ // Note: We are using hashchange instead of popstate
        if ("#!" + currentState != location.hash) {
            if ((undefined != location.hash) && ("" != location.hash)) {
                showContent(location.hash);      
            } else {
                showContent($(settings.listItem + ':first').attr("name"));
            }
        }
    });

    
   function showContent(hash){
       var hashReplace = hash.replace("#!", ""); 
       $(settings.listItem).removeClass('selected');
       $("#" + hashReplace + "-tab").addClass('selected');
       var newTitle = $("#" + hashReplace + "-tab a:first").attr("title");
       var div = $("#" + hashReplace + "-content").html();
       document.title = newTitle;
       currentState = hashReplace;
       $(settings.contentItem).replaceWith(div);
       // $.ajax({
       //   url: $.url().attr("base") + $.url().attr("directory") + $.url().attr("file") + "?_escaped_fragment_=" +  hashReplace,
       //   dataType: "html",
       //   success: function(html) {
       //             var div = $(settings.contentItem, $(html)).addClass('done');
       //             $(settings.contentItem).replaceWith(div);          
       //             //initPage();
       //   }
       // });    
   }
    

};

var featureStickyNavInit = function(){
  var scrollBottom = $(window).scrollTop() + $(window).height();
    var documentHeight = $(document).height();
    var footerHeight = $('#footer').height();

    stickNavigation();

    if ((documentHeight - scrollBottom) < footerHeight) {
        $('#featureNav').hide();
        $(window).scroll();
        setTimeout("$('#featureNav').fadeIn('fast')", 100);
    }

    // setTimeout("$(window).trigger('scroll')", 100);


    $('.tour-nav a').click(function() {
        $('.tour-nav a').closest('li').removeClass('selected');
        $(this).closest('li').addClass('selected');
    });

    $(window).resize(function() {
        stickNavigation();
    });
    // $('.tour-nav').localScroll({'hash':true});
  
  function stickNavigation() {

      var headerHeight = $('#featureNav').offset().top;
      var navPosition = $('#featureNav').offset().top;
      var navOffsetPosition = $('.stickyNav').offset().top;
    var footerPosition = $("#endOfTour").offset().top;
      var lastItemPosition = $('.navClone').offset().top + ($('.navClone').outerHeight() + 30);

      var windowHeight = $(window).height();
      var navHeight = $('#featureNav').height();
    var distanceFromBottom = windowHeight- navHeight;

      $('#featureList').css('position', 'relative');

      $(".navClone").css({
          'height': $('#featureNav').height(),
          'width': $('#featureNav').outerWidth(),
          'visibility' : 'hidden',
          'margin-left' : -$('#featureNav').outerWidth(),
          'float' : 'left'
      });


      if ($('#featureNav').height() > $(window).height()) {
          var navHeightOffset = $('.stickyNav').height() - $(window).height();
      } else {
          var navHeightOffset = 0;
      }
                 
    function setNavOnScroll() {
      navPosition = $('#featureNav').offset().top;
        navOffsetPosition = $('.stickyNav').offset().top;
        footerPosition = $("#endOfTour").offset().top;
        lastItemPosition = parseInt($('.navClone').offset().top + ($('.navClone').height()));
        originalItemPosition = parseInt($('#featureNav').offset().top + ($('#featureNav').height()));

        setNavPosition({
            headerHeight : headerHeight,
            navPosition : navPosition,
            navOffsetPosition : navOffsetPosition,
            footerPosition : footerPosition,
            lastItemPosition : lastItemPosition,
            navHeightOffset : navHeightOffset,
        distanceFromBottom : distanceFromBottom 
        });

        // console.log($(document).scrollTop());
        // console.log($('#spaces').offset().top);

    if (($(window).width() < 960)) {
      if ($('.stickyNav').hasClass('stuck')) {
        $('#featureNav li.selected').css({'zIndex' : '0'});
      }
    } else {
      $('#featureNav li.selected').css({'z-index' : '10'});
    }
        $('.tabbed-section').each(function(i, a) {
            var currentPos = $(a).offset().top - 20;
            if ($(document).scrollTop() > currentPos) {
                var elName = "." + $(a).attr('id');
                $('.tour-nav li').removeClass('selected');
                $(elName).closest('li').addClass('selected');
            }
        });
    };


    // $(window).scroll( $.debounce( 20, setNavOnScroll ) );                         
    $(window).scroll(function(){
           setNavOnScroll(); 
    }); 

  $('.backToTop').click(function(e) {
    e.preventDefault();
        $('.tour-nav li').removeClass('selected');
        $('.tour-nav li:first').addClass('selected');
    $(window).scrollTo($("#featureList"),0);
    });


  }

  function setNavPosition(elValues) {



      if (elValues.lastItemPosition <= elValues.footerPosition - elValues.distanceFromBottom) {
          if ($(document).scrollTop() >= (elValues.headerHeight + elValues.navHeightOffset)) {
              $('.stickyNav').css({'position' : 'fixed', 'top' : - elValues.navHeightOffset}).addClass('stuck');
          } else {
              $('.stickyNav').css({'position' : 'static'}).removeClass('stuck');

          }
      } else {
          $('#featureNav').css({'position' : 'absolute', 'bottom' : elValues.distanceFromBottom, 'top' : 'auto'});
      }

  }  




};

var lightBoxVideo = function(videoNum){
  
var customHeight = $('.videoLB').data('height');
var customWidth = $('.videoLB').data('width');

var settings = {
        iframe: true,
        opacity: 0.7,
      transition: 'none',
      innerWidth:640, 
      innerHeight:390  
}

if (customHeight > 0) {
  settings.innerHeight = customHeight;
  settings.innerWidth = customWidth;
}
  
  $('.videoLB').colorbox({
        iframe: settings.iframe,
        opacity: settings.opacity,
    transition: settings.transition,
  innerWidth:settings.innerWidth, 
  innerHeight:settings.innerHeight,             
maxWidth:"100%",
maxHeight:"100%",
        onOpen: function() {
      // $(videoEmbedDiv).show();
            $('body').addClass('videoOpen');
        },

    onCleanup: function() {
      // $(videoEmbedDiv).hide();
    },
    
  onClosed: function() {
    $('body').removeClass('videoOpen');
        }
    });
};

var removeCurveFirstItem = function() {
  if ($('#featureNav li:first').hasClass('selected')) {
        $('#tourContent').addClass('navTop');
    } else {
        $('#tourContent').removeClass('navTop');
    }
}; 

var orderJobList = function(listToEffect){
  //loads slow, would be cool to detect when the css has been set and then fade in
  if (!listToEffect) {
    listToEffect = ".all";
  }                                                    

    var numberOfItems = $(listToEffect +' .jobListings li').length;
    var numberPerColumn = Math.ceil(numberOfItems / 3);
    var count = 0;
    var group = 1;
        
    //if there is less than 10 items dont split them up between cols
        if (numberOfItems < 10) {
      var numberPerColumn = 9;
    }    
    

    $(listToEffect +' .jobListings li').each(function(i, item){
      if (count == numberPerColumn) {
        count=0;
        group++;
      }

      // $(this).addClass('listCol-' + group);
      
      // use the group number to append it to the ul with that number
      $(this).appendTo(listToEffect + ' .group-'+group);
      
      
      count++;

    });  
    
    $(listToEffect +' .heading').each(function(){   
      if ($(this).next().length == 0 || $(this).next().hasClass('heading')) {
          $(this).css('display', 'none');
      }

    });
};  

//This is currently not being used
// var googleSearch = function() {
// 
//     $( "#searchInput" ).autocomplete({
//       source: function( request, response ) {
//         $.ajax({
//           url: "http://clients1.google.com/complete/search",
//           dataType: "jsonp",
//           data: {
//             ds: "cse",
//             client: "partner",
//             source: "gcsc",
//             partnerid: "004158415757119463896:yiaveaifukw",
//             hl: "en",
//             q: request.term
//           },
//           success: function( data ) {
//             response( $.map( data[1], function( item ) {
//               return {
//                 label: item[0],
//                 value: item[0]
//               }
//             }));
//           }
//         });
//       },
//       minLength: 2,
//       select: function( event, ui ) {
//         log( ui.item ?
//           "Selected: " + ui.item.label :
//           "Nothing selected, input was " + this.value);
//         $("#searchForm").submit();
//       },
//       open: function() {
//         $( this ).removeClass( "ui-corner-all" ).addClass( "ui-corner-top" );
//       },
//       close: function() {
//         $( this ).removeClass( "ui-corner-top" ).addClass( "ui-corner-all" );
//       }
//      });  
//      
//      
// }

var confluencePricing = function() {

  var downloadDefPrice = $('#download .defPrice').html().replace(',','');
  var onDemandDefPrice = $('#ondemand .defPrice').html().replace(',','');

  downloadDefPrice = parseInt(downloadDefPrice);
  onDemandDefPrice = parseInt(onDemandDefPrice);


  var defPrice = parseInt($('.defPrice').html()); 

    function addRemovePrice(options){

      var prices = $('.unitPrice:visible');

      prices.each(function(){

        var unitPrice = parseInt($(this).html().replace(',', '').replace('.', ''));
        var addOnProductPrice= parseInt($(this).siblings(options.productCheckBoxSelector).html().replace(',', '').replace('.', ''));      
        
        // start comma function
        function commaize(n) {
            n=n+'';
            return (Math.abs(n) <= 999) ? n : commaize(n.slice(0,-3))+','+n.slice(-3);
        }
        //end


        if (options.addOrRemove == "add") {
          var newPrice = unitPrice + addOnProductPrice;
          var commaPrice = commaize(newPrice);
          $(this).html(commaPrice); 
        } else {
          var newPrice = unitPrice - addOnProductPrice;   
          var commaPrice = commaize(newPrice);        
          $(this).html(commaPrice);
        }                                         
      });

    }      

    function changeStarter(addRemove) {
      var downloadStarterPriceNode = $('#download .defPrice').html().replace(',','');
      var onDemandStarterPriceNode = $('#ondemand .defPrice').html().replace(',','');

      if (addRemove == 'add') {
        downloadStarterPriceNode = parseInt(downloadStarterPriceNode) + downloadDefPrice;
        onDemandStarterPriceNode = parseInt(onDemandStarterPriceNode) + onDemandDefPrice;
      } else {                                                        
        downloadStarterPriceNode = parseInt(downloadStarterPriceNode) - downloadDefPrice;
        onDemandStarterPriceNode = parseInt(onDemandStarterPriceNode) - onDemandDefPrice;
      }
      
      $('#download .defPrice').html(downloadStarterPriceNode);
      $('#ondemand .defPrice').html(onDemandStarterPriceNode);    
    }

    History.Adapter.bind(window,'statechange',function(){
      $('.defPrice').html(defPrice);
       $('#download .defPrice').html(downloadDefPrice);
       $('#ondemand .defPrice').html(onDemandDefPrice);
      
      $(':checkbox:checked').each(function(){
        if ($(this).is(':visible')) {
           changeStarter('add');
        }

      });
    }); 

    $(':checkbox:checked').removeAttr('checked'); 

    $('.tcCheckbox').change(function(){
    if ($(this).is(':checked')) {  
      addRemovePrice({
        productCheckBoxSelector : '.teamCalPrice',
        addOrRemove : 'add'
      });
      
      changeStarter('add');

    } else {                        
      addRemovePrice({
        productCheckBoxSelector : '.teamCalPrice',
        addOrRemove : 'remove'
      });
      changeStarter('remove');      
    }
    });

    $('.sharepointCheckbox').change(function(){

    if ($(this).is(':checked')) {  
      addRemovePrice({
        productCheckBoxSelector : '.sharepointPrice',
        addOrRemove : 'add'
      });                                 
      changeStarter('add');      

    } else {                        
      addRemovePrice({
        productCheckBoxSelector : '.sharepointPrice',
        addOrRemove : 'remove'
      }); 
      changeStarter('remove');      
    }
    });

   
   
  
};
  
var jiraPricing = function() {
                                                        


var downloadDefPrice = $('#download .defPrice').html().replace(',','');
var onDemandDefPrice = $('#ondemand .defPrice').html().replace(',','');

downloadDefPrice = parseInt(downloadDefPrice);
onDemandDefPrice = parseInt(onDemandDefPrice);


var defPrice = parseInt($('.defPrice').html()); 
  
  function addRemovePrice(options){

    var prices = $('.unitPrice:visible');
    
    prices.each(function(){
      var unitPrice = parseInt($(this).html().replace(',', '').replace('.', ''));
      var addOnProductPrice= parseInt($(this).siblings(options.productCheckBoxSelector).html().replace(',', '').replace('.', ''));      
      
      // start comma function
      function commaize(n) {
          n=n+'';
          return (Math.abs(n) <= 999) ? n : commaize(n.slice(0,-3))+','+n.slice(-3);
      }
      //end
      
      if (options.addOrRemove == "add") {
        var newPrice = unitPrice + addOnProductPrice;
        var commaPrice = commaize(newPrice);
        $(this).html(commaPrice); 
      } else {
        var newPrice = unitPrice - addOnProductPrice;   
        var commaPrice = commaize(newPrice);        
        $(this).html(commaPrice);
      }                                         
    });
    
  }      
  
  function changeStarter(addRemove) {
    var downloadStarterPriceNode = $('#download .defPrice').html().replace(',','');
    var onDemandStarterPriceNode = $('#ondemand .defPrice').html().replace(',','');
        
    if (addRemove == 'add') {

            
      downloadStarterPriceNode = parseInt(downloadStarterPriceNode) + downloadDefPrice;
      onDemandStarterPriceNode = parseInt(onDemandStarterPriceNode) + onDemandDefPrice;
      
    } else {                                                        
      
            
      downloadStarterPriceNode = parseInt(downloadStarterPriceNode) - downloadDefPrice;
      onDemandStarterPriceNode = parseInt(onDemandStarterPriceNode) - onDemandDefPrice;
    }

    
    $('#download .defPrice').html(downloadStarterPriceNode);
    $('#ondemand .defPrice').html(onDemandStarterPriceNode);    
  }
  
  History.Adapter.bind(window,'statechange',function(){
    $('.defPrice').html(defPrice);
     $('#download .defPrice').html(downloadDefPrice);
     $('#ondemand .defPrice').html(onDemandDefPrice);

    
    $(':checkbox:checked').each(function(){
      if ($(this).is(':visible')) {
         changeStarter('add');
      }

    });
  }); 

  $(':checkbox:checked').removeAttr('checked'); 

    $('.greenhopperCheckbox').change(function(){
    if ($(this).is(':checked')) {  
      addRemovePrice({
        productCheckBoxSelector : '.greenhopperPrice',
        addOrRemove : 'add'
      });

            changeStarter('add');

    } else {                        
      addRemovePrice({
        productCheckBoxSelector : '.greenhopperPrice',
        addOrRemove : 'remove'
      });
      
      changeStarter('remove');
            
    }
    });
                                          
    $('.bonfireCheckbox').change(function(){
         // checkBonfire();
    if ($(this).is(':checked')) {  
      addRemovePrice({
        productCheckBoxSelector : '.bonfirePrice',
        addOrRemove : 'add'
      });
      changeStarter('add');

    } else {                        
      addRemovePrice({
        productCheckBoxSelector : '.bonfirePrice',
        addOrRemove : 'remove'
      });
      changeStarter('remove');
    }
    });

  $('.gliffyCheckbox').change(function(){ 
  if ($(this).is(':checked')) {  
    addRemovePrice({
      productCheckBoxSelector : '.gliffyPrice',
      addOrRemove : 'add'
    });
    changeStarter('add');
    
  } else {                        
    addRemovePrice({
      productCheckBoxSelector : '.gliffyPrice',
      addOrRemove : 'remove'
    });
     changeStarter('remove');
  } 
         
  }); 
  
 
  
}; 

var sharepointConnectorPricing = function(){
   var defPrice = $('.defPrice').html(); 
    
    function showConfluenceOnly(){
      $('.confluencePrice').show();      
      $('.unitPrice').hide();
      $('.defPrice').html(defPrice * 2);      
    }   
    
    
    function showDefault(){
      $('.unitPrice').show();      
      $('.confluencePrice').hide();
      $('.defPrice').html(defPrice);      
    }

        $(':checkbox:checked').removeAttr('checked');  
    
    $('.confluenceCheckbox').change(function(){
      checkConfluenceCheckbox();
    });
                                                          
    $("#tabsNav a").click(function(){
      $(':checkbox:checked').removeAttr('checked');  
      checkConfluenceCheckbox();     
    });
    
      
    function checkConfluenceCheckbox(){
      if ($('.confluenceCheckbox').is(':checked')){
         showConfluenceOnly();   
      } else {
        showDefault();         
      }
    }
};

var bonfirePricing = function(){
  var starterPrice = 10;
  
  function addRemovePrice(options){

    var prices = $('.unitPrice:visible');
    
    prices.each(function(){
      var unitPrice = parseInt($(this).html().replace(',', '').replace('.', ''));
      var addOnProductPrice= parseInt($(this).siblings(options.productCheckBoxSelector).html().replace(',', '').replace('.', ''));      

      // start comma function
      function commaize(n) {
          n=n+'';
          return (Math.abs(n) <= 999) ? n : commaize(n.slice(0,-3))+','+n.slice(-3);
      }
      //end
      
      if (options.addOrRemove == "add") {
        var newPrice = unitPrice + addOnProductPrice;
        var commaPrice = commaize(newPrice);
        $(this).html(commaPrice); 
      } else {
        var newPrice = unitPrice - addOnProductPrice;   
        var commaPrice = commaize(newPrice);        
        $(this).html(commaPrice);
      }                                         
    });
    
  }
  function changeStarter(addRemove) {
    
    if (addRemove == 'add') {
      var starterPriceNode = parseInt($('.defPrice').html()) + starterPrice;
    } else {
      var starterPriceNode = parseInt($('.defPrice').html()) - starterPrice;      
    }

    
    $('.defPrice').html(starterPriceNode);
  }
      History.Adapter.bind(window,'statechange',function(){
    $('.defPrice').html('10');
    $(':checkbox:checked').each(function(){
      if ($(this).is(':visible')) {
         changeStarter('add');
      }

    });
  });
  $(':checkbox:checked').removeAttr('checked'); 

    $('.jiraCheckbox').change(function(){
    if ($(this).is(':checked')) {  
      addRemovePrice({
        productCheckBoxSelector : '.jiraPrice',
        addOrRemove : 'add'
      }); 
      changeStarter('add');

    } else {                        
      addRemovePrice({
        productCheckBoxSelector : '.jiraPrice',
        addOrRemove : 'remove'
      });
      changeStarter('remove');
    }
    });

};
  
var fisheyePricing = function(){

  var starterPrice = 10;
  
  function addRemovePrice(options){

    var prices = $('.unitPrice:visible');
    
    prices.each(function(){
      var unitPrice = parseInt($(this).html().replace(',', '').replace('.', ''));
      var addOnProductPrice= parseInt($(this).closest('li').find(options.productCheckBoxSelector).html().replace(',', '').replace('.', ''));      

      // start comma function
      function commaize(n) {
          n=n+'';
          return (Math.abs(n) <= 999) ? n : commaize(n.slice(0,-3))+','+n.slice(-3);
      }
      //end
      
      if (options.addOrRemove == "add") {
        var newPrice = unitPrice + addOnProductPrice;
        var commaPrice = commaize(newPrice);
        $(this).html(commaPrice); 
      } else {
        var newPrice = unitPrice - addOnProductPrice;   
        var commaPrice = commaize(newPrice);        
        $(this).html(commaPrice);
      }                                         
    });
    
  }
  function changeStarter(addRemove) {
    
    if (addRemove == 'add') {
      var starterPriceNode = parseInt($('.defPrice').html()) + starterPrice;
    } else {
      var starterPriceNode = parseInt($('.defPrice').html()) - starterPrice;      
    }

    
    $('.defPrice').html(starterPriceNode);
  }
      History.Adapter.bind(window,'statechange',function(){
    $('.defPrice').html('10');
    $(':checkbox:checked').each(function(){
      if ($(this).is(':visible')) {
         changeStarter('add');
      }

    });
  });
  $(':checkbox:checked').removeAttr('checked'); 

    $('.jiraCheckbox').change(function(){
    if ($(this).is(':checked')) {  
      addRemovePrice({
        productCheckBoxSelector : '.jiraPrice',
        addOrRemove : 'add'
      }); 
      changeStarter('add');

    } else {                        
      addRemovePrice({
        productCheckBoxSelector : '.jiraPrice',
        addOrRemove : 'remove'
      });
      changeStarter('remove');
    }
    });

};

var bambooPricing = function(){

  var starterPrice = 10;
  
  function addRemovePrice(options){

    var prices = $('.unitPrice:visible');
    
    prices.each(function(){
      var unitPrice = parseInt($(this).html().replace(',', '').replace('.', ''));
      var addOnProductPrice= parseInt($(this).closest('li').find(options.productCheckBoxSelector).html().replace(',', '').replace('.', ''));      

      if (options.addOrRemove == "add") {
        $(this).html(unitPrice + addOnProductPrice);
      } else {                                  
        $(this).html(unitPrice - addOnProductPrice);
      }                                         
    });
    
  }
  function changeStarter(addRemove) {
    
    if (addRemove == 'add') {
      var starterPriceNode = parseInt($('.defPrice').html()) + starterPrice;
    } else {
      var starterPriceNode = parseInt($('.defPrice').html()) - starterPrice;      
    }

    
    $('.defPrice').html(starterPriceNode);
  }
      History.Adapter.bind(window,'statechange',function(){
    $('.defPrice').html('10');
    $(':checkbox:checked').each(function(){
      if ($(this).is(':visible')) {
         changeStarter('add');
      }

    });
  });
  $(':checkbox:checked').removeAttr('checked'); 

    $('.fisheyeCheckbox').change(function(){
    if ($(this).is(':checked')) {  
      addRemovePrice({
        productCheckBoxSelector : '.fisheyePrice',
        addOrRemove : 'add'
      }); 
      changeStarter('add');

    } else {                        
      addRemovePrice({
        productCheckBoxSelector : '.fisheyePrice',
        addOrRemove : 'remove'
      });
      changeStarter('remove');
    }
    });
       
    $('.crucibleCheckbox').change(function(){
    if ($(this).is(':checked')) {  
      addRemovePrice({
        productCheckBoxSelector : '.cruciblePrice',
        addOrRemove : 'add'
      }); 
      changeStarter('add');

    } else {                        
      addRemovePrice({
        productCheckBoxSelector : '.cruciblePrice',
        addOrRemove : 'remove'
      });
      changeStarter('remove');
    }
    });

};


var greenHopperPricing = function(){
  var starterPrice = 10;      
  var defPrice = $('.defPrice').html();

  function addRemovePrice(options){
    var prices = $('.unitPrice:visible');

    prices.each(function(){
      var unitPrice = parseInt($(this).html().replace(',', '').replace('.', ''));
      var addOnProductPrice= parseInt($(this).siblings(options.productCheckBoxSelector).html().replace(',', '').replace('.', ''));      

      // start comma function
      function commaize(n) {
        n=n+'';
        return (Math.abs(n) <= 999) ? n : commaize(n.slice(0,-3))+','+n.slice(-3);
      }
      //end

      if (options.addOrRemove == "add") {
        var newPrice = unitPrice + addOnProductPrice;
        var commaPrice = commaize(newPrice);
        $(this).html(commaPrice); 
      } else {
        var newPrice = unitPrice - addOnProductPrice;   
        var commaPrice = commaize(newPrice);        
        $(this).html(commaPrice);
      }                                         
    });
  }
  function changeStarter(addRemove) {
    if (addRemove == 'add') {
      var starterPriceNode = parseInt($('.defPrice').html()) + starterPrice;
    } else {
      var starterPriceNode = parseInt($('.defPrice').html()) - starterPrice;      
    }
    $('.defPrice').html(starterPriceNode);
  }
  History.Adapter.bind(window,'statechange',function(){
    $('.defPrice').html('10');
    $(':checkbox:checked').each(function(){
      if ($(this).is(':visible')) {
        changeStarter('add');
      }
    });
  });
  $(':checkbox:checked').removeAttr('checked'); 

  $('.jiraCheckbox').change(function(){
    if ($(this).is(':checked')) {  
      addRemovePrice({
      productCheckBoxSelector : '.jiraPrice',
      addOrRemove : 'add'
      });
      changeStarter('add');
    } else {                        
      addRemovePrice({
      productCheckBoxSelector : '.jiraPrice',
      addOrRemove : 'remove'
      }); 
      changeStarter('remove');
    }
  });
  $(':checkbox:checked').removeAttr('checked');
};

var stashPricing = function(){
  var starterPrice = 10;      
var defPrice = $('.defPrice').html();

    function addRemovePrice(options){
      var prices = $('.unitPrice:visible');
      
      prices.each(function(){
        var unitPrice = parseInt($(this).html().replace(',', '').replace('.', ''));
        var addOnProductPrice= parseInt($(this).siblings(options.productCheckBoxSelector).html().replace(',', '').replace('.', ''));      
        var afterSalePrice = $(this).closest('li').find('.afterSalePrice').html();
        if (afterSalePrice != null) {
          afterSalePrice = parseInt(afterSalePrice.replace(',', '').replace('.', ''));          
        }

        if (options.addOrRemove == "add") {
          $(this).html(unitPrice + addOnProductPrice);
          $(this).closest('li').find('.afterSalePrice').html(afterSalePrice + addOnProductPrice);
        } else {                                  
          $(this).html(unitPrice - addOnProductPrice);
          $(this).closest('li').find('.afterSalePrice').html(afterSalePrice - addOnProductPrice);          
        }                                         
      });
      

    }
    function changeStarter(addRemove) {

      if (addRemove == 'add') {
        var starterPriceNode = parseInt($('.defPrice').html()) + starterPrice;
      } else {
        var starterPriceNode = parseInt($('.defPrice').html()) - starterPrice;      
      }


      $('.defPrice').html(starterPriceNode);
    }
            History.Adapter.bind(window,'statechange',function(){
      $('.defPrice').html('10');
      $(':checkbox:checked').each(function(){
        if ($(this).is(':visible')) {
           changeStarter('add');
        }

      });
    });
    $(':checkbox:checked').removeAttr('checked'); 

      $('.jiraCheckbox').change(function(){
      if ($(this).is(':checked')) {  
        addRemovePrice({
          productCheckBoxSelector : '.jiraPrice',
          addOrRemove : 'add'
        });
        changeStarter('add');

      } else {                        
        addRemovePrice({
          productCheckBoxSelector : '.jiraPrice',
          addOrRemove : 'remove'
        }); 
        changeStarter('remove');
      }
      });


      $(':checkbox:checked').removeAttr('checked');
};



/*var teamCalendarsPricing = function(){
  var defPrice = $('.defPrice').html();

  function showConfluenceOnly(){
    $('.confluencePrice').show();      
    $('.unitPrice').hide();
    $('.defPrice').html(defPrice * 2);      
  }


  function showDefault(){
    $('.unitPrice').show();      
    $('.confluencePrice').hide();
    $('.defPrice').html(defPrice);       
  }

  $(':checkbox:checked').removeAttr('checked');  

  $('.confluenceCheckbox').change(function(){
   checkConfluenceCheckbox();
  });

  $("#tabsNav a").click(function(){
    $(':checkbox:checked').removeAttr('checked');  
    checkConfluenceCheckbox();     
  });


  function checkConfluenceCheckbox(){
    if ($('.confluenceCheckbox').is(':checked')){
      showConfluenceOnly();   
    } else {
      showDefault();         
    }
  }
};*/

var teamCalendarsPricing = function(){
  var starterPrice = 10;      
  var defPrice = $('.defPrice').html();

  function addRemovePrice(options){
    var prices = $('.unitPrice:visible');

    prices.each(function(){
      var unitPrice = parseInt($(this).html().replace(',', '').replace('.', ''));
      var addOnProductPrice= parseInt($(this).siblings(options.productCheckBoxSelector).html().replace(',', '').replace('.', ''));      

      // start comma function
      function commaize(n) {
        n=n+'';
        return (Math.abs(n) <= 999) ? n : commaize(n.slice(0,-3))+','+n.slice(-3);
      }
      //end

      if (options.addOrRemove == "add") {
        var newPrice = unitPrice + addOnProductPrice;
        var commaPrice = commaize(newPrice);
        $(this).html(commaPrice); 
      } else {
        var newPrice = unitPrice - addOnProductPrice;   
        var commaPrice = commaize(newPrice);        
        $(this).html(commaPrice);
      }                                         
    });
  }
  function changeStarter(addRemove) {
    if (addRemove == 'add') {
      var starterPriceNode = parseInt($('.defPrice').html()) + starterPrice;
    } else {
      var starterPriceNode = parseInt($('.defPrice').html()) - starterPrice;      
    }
    $('.defPrice').html(starterPriceNode);
  }
  History.Adapter.bind(window,'statechange',function(){
    $('.defPrice').html('10');
    $(':checkbox:checked').each(function(){
      if ($(this).is(':visible')) {
        changeStarter('add');
      }
    });
  });
  $(':checkbox:checked').removeAttr('checked'); 

  $('.confluenceCheckbox').change(function(){
    if ($(this).is(':checked')) {  
      addRemovePrice({
      productCheckBoxSelector : '.confluencePrice',
      addOrRemove : 'add'
      });
      changeStarter('add');
    } else {                        
      addRemovePrice({
      productCheckBoxSelector : '.confluencePrice',
      addOrRemove : 'remove'
      }); 
      changeStarter('remove');
    }
  });
  $(':checkbox:checked').removeAttr('checked');
};

var sharepointConnectorPricing = function(){
  var starterPrice = 10;      
  var defPrice = $('.defPrice').html();

  function addRemovePrice(options){
    var prices = $('.unitPrice:visible');

    prices.each(function(){
      var unitPrice = parseInt($(this).html().replace(',', '').replace('.', ''));
      var addOnProductPrice= parseInt($(this).siblings(options.productCheckBoxSelector).html().replace(',', '').replace('.', ''));      

      // start comma function
      function commaize(n) {
        n=n+'';
        return (Math.abs(n) <= 999) ? n : commaize(n.slice(0,-3))+','+n.slice(-3);
      }
      //end

      if (options.addOrRemove == "add") {
        var newPrice = unitPrice + addOnProductPrice;
        var commaPrice = commaize(newPrice);
        $(this).html(commaPrice); 
      } else {
        var newPrice = unitPrice - addOnProductPrice;   
        var commaPrice = commaize(newPrice);        
        $(this).html(commaPrice);
      }                                         
    });
  }
  function changeStarter(addRemove) {
    if (addRemove == 'add') {
      var starterPriceNode = parseInt($('.defPrice').html()) + starterPrice;
    } else {
      var starterPriceNode = parseInt($('.defPrice').html()) - starterPrice;      
    }
    $('.defPrice').html(starterPriceNode);
  }
  History.Adapter.bind(window,'statechange',function(){
    $('.defPrice').html('10');
    $(':checkbox:checked').each(function(){
      if ($(this).is(':visible')) {
        changeStarter('add');
      }
    });
  });
  $(':checkbox:checked').removeAttr('checked'); 

  $('.confluenceCheckbox').change(function(){
    if ($(this).is(':checked')) {  
      addRemovePrice({
      productCheckBoxSelector : '.confluencePrice',
      addOrRemove : 'add'
      });
      changeStarter('add');
    } else {                        
      addRemovePrice({
      productCheckBoxSelector : '.confluencePrice',
      addOrRemove : 'remove'
      }); 
      changeStarter('remove');
    }
  });
  $(':checkbox:checked').removeAttr('checked');
};

var perkSlider = function(settings){
  
  if (!settings) {
      var settings = {
      slideContainerClass:'.slideshow',
      keyboardControl: true
    } 
  }
          $(window).resize(function(){

        if ($(window).width() < 960) {
          $('body').css('overflow-x', 'scroll');
        } else {
          $('body').css('overflow-x', 'hidden');
        }

      });


      var slideEl = $(settings.slideContainerClass).find('li');
      var slideWidth = slideEl.outerWidth();
      var totalWidth = (slideEl.length * slideEl.outerWidth());

      var nextSlidePos = -slideWidth; 
      var prevSlidePos = 0;
      var slidePosition = 0;

  /*
    TODO                                             

    try making this better using this technique
    var wait = setInterval(function() {
      if( !$("#element1, #element2").is(":animated") ) {
        clearInterval(wait);
        // This piece of code will be executed
        // after element1 and element2 is complete.
      }
    }, 200);

  */  


      $('.slideshow').css('width', totalWidth);

      if (parseInt($('.slideshow').css('marginLeft'), 10) > 0) {
        $('.previous').fadeIn('fast');
      }

      $('.next').fadeIn('fast');

      $('.next').live('click', function(e){
        e.preventDefault();
        slidePosition =  parseInt($('.slideshow').css('marginLeft'), 10);
        nextSlidePos = slidePosition + (-slideWidth);

        if (!$('.slideshow').is(':animated') && (-parseInt($('.slideshow').css('marginLeft'), 10) + slideWidth) !== totalWidth) {
        $('.slideshow').animate({marginLeft: nextSlidePos}, 500, function(){
          slidePosition =  parseInt($('.slideshow').css('marginLeft'), 10);

          nextSlidePos = slidePosition + (-slideWidth);

          if (-parseInt($('.slideshow').css('marginLeft'), 10) > 0) {
            $('.previous').fadeIn('fast');
          } 

          if ((-parseInt($('.slideshow').css('marginLeft'), 10) + slideWidth) === totalWidth) {
            $('.next').fadeOut('fast');
          }


        });
        }
      });     

      $('.previous').live('click', function(e){
        e.preventDefault();
        slidePosition =  parseInt($('.slideshow').css('marginLeft'), 10);
        prevSlidePos = slidePosition + slideWidth;                     

        if (!$('.slideshow').is(':animated') && (-parseInt($('.slideshow').css('marginLeft'), 10) > 0)) {
          $('.slideshow').animate({marginLeft: prevSlidePos}, 500, function(){
            slidePosition =  parseInt($('.slideshow').css('marginLeft'), 10);
            prevSlidePos = slidePosition + slideWidth;
            if (-parseInt($('.slideshow').css('marginLeft'), 10) > 0) {

              $('.previous').fadeIn('fast');
            } else {
              $('.previous').fadeOut('fast');
            } 

            if ((-parseInt($('.slideshow').css('marginLeft'), 10) + slideWidth) < totalWidth) {
              $('.next').fadeIn('fast');
            } 
          });
        }

      });    
      
  if (settings.keyboardControl == true) {
    $('body').keyup(function(e) {
        if (e.keyCode == '39') {
          e.preventDefault();
          if ($('.next').is(':visible')) {
            $('.next a').click();
          }

         }
      });

      $('body').keyup(function(e) {
        if (e.keyCode == '37') {
           e.preventDefault();
        if ($('.previous').is(':visible')) {
          $('.previous a').click();
        }
         }
      });
  }


};

                                     
$(function(){
    if (!Array.prototype.indexOf) {
      Array.prototype.indexOf = function (searchElement /*, fromIndex */ ) {
          "use strict";
          if (this == null) {
              throw new TypeError();
          }
          var t = Object(this);
          var len = t.length >>> 0;
          if (len === 0) {
              return -1;
          }
          var n = 0;
          if (arguments.length > 1) {
              n = Number(arguments[1]);
              if (n != n) { // shortcut for verifying if it's NaN
                  n = 0;
              } else if (n != 0 && n != Infinity && n != -Infinity) {
                  n = (n > 0 || -1) * Math.floor(Math.abs(n));
              }
          }
          if (n >= len) {
              return -1;
          }
          var k = n >= 0 ? n : Math.max(len - Math.abs(n), 0);
          for (; k < len; k++) {
              if (k in t && t[k] === searchElement) {
                  return k;
              }
          }
          return -1;
      }
  }
  $('.fancyIntro').each(function(){
    $(this).find('p:first').fancyletter({commonClass: 'highlite'});
  });
  $('#globalNav .languageSelect').hover(function() {
    $('.languages').slideDown('fast');
  }, function() {
    $('.languages').slideUp('fast');
  });  

  //   setSameLiHeight('.brandsList li');
  //   $('.tourLink').click(function() {
  //   
  // // setTimeout(function() {setSameLiHeight('.brandsList li');}, 0);
  //     
  //   });           

});


var setSameLiHeight = function (liItem) {
    var highestLi = 0;
    
    $(liItem).each(function(){
      liHeight = $(this).height();
    if (liHeight > highestLi) {
      highestLi = liHeight;
    
    }
    }); 

    $(liItem).css('height', highestLi).fadeTo(100,1);
}

var setSameRowHeight = function (liItem) { 
 var currentTallest = 0,
   currentRowStart = 0,
   rowDivs = new Array(),
   $el,
   topPosition = 0;
 
  $(liItem).each(function() {
 
    $el = $(this);
    topPostion = $el.position().top;
 
    if (currentRowStart != topPostion) {
 
    // we just came to a new row.  Set all the heights on the completed row
    for (currentDiv = 0 ; currentDiv < rowDivs.length ; currentDiv++) {
      rowDivs[currentDiv].height(currentTallest).fadeTo(100,1);
    }
 
    // set the variables for the new row
    rowDivs.length = 0; // empty the array
    currentRowStart = topPostion;
    currentTallest = $el.height();
    rowDivs.push($el);
 
    } else {
 
    // another div on the current row.  Add it to the list and check if it's taller
    rowDivs.push($el);
    currentTallest = (currentTallest < $el.height()) ? ($el.height()) : (currentTallest);
 
   }
 
   // do the last row
    for (currentDiv = 0 ; currentDiv < rowDivs.length ; currentDiv++) {
    rowDivs[currentDiv].height(currentTallest).fadeTo(100,1);
    }
 
  });    
       
}

var initVideos = function(){

    if ($('.launchVideo').length > 1) {

      $('.playButton').each(function(i){
        $(this).addClass('videoLB-' + i);
      });

      $('.playButton').click(function() {                
            var videoNumber = $(this).prev().attr('rel');
            lightBoxVideo(videoNumber);        
        });

      $('.launchVideo').each(function(i){
      if (!$(this).parent().hasClass('caption')) {
      $(this).attr('rel', i);
      }
        $(this).addClass('videoLB-' + i);
      });
                
      $('.hiddenEmbed').each(function(i){
        $(this).addClass('videoEmbed-' + i);
      });    

        $('.launchVideo').click(function() {
            var videoNumber = $(this).attr('rel');
            lightBoxVideo(videoNumber);
        });
    } else {
          $('.hiddenEmbed').addClass('videoEmbed');
      lightBoxVideo();
    }
  
  // $('.videoWrapper').find();
  
} 

var initExpertSearch = function(){
      // Clear the form on load, timeout is for IE
    setTimeout(function(){
      document.getElementById("searchPartner").reset();
    }, 5);

    $('.resultsMessage').hide();        
    function expertSearchResults(){
      /*
        TODO  could this have easily been magnolia ctx?
      */

      var expertiseSearched = $.url().param('expertskills');
      var productsSearched = $.url().param('products');
      var isPremier = $.url().param('isPremier');
      var expertLocation = $.url().param('expertlocation');
            

            
            

        if (!expertiseSearched) {
          expertiseSearched = '*';
        }                           

        if (!productsSearched) {
          productsSearched = "*";
        }
                                               
        if (!expertLocation) {
          expertLocation = "";
        }
        
        
        if (!isPremier) {
          isPremier = "*";
        } else {                       
          $('#ispremier_box').attr('checked','checked'); 
        }

                var expertiseSearchedArray = expertiseSearched.split(',');

                for (var i=0; i < expertiseSearchedArray.length; i++) {
                    if (expertiseSearchedArray[i] == "ManagedHosting") {
                        expertiseSearchedArray[i] = "Managed Hosting";
                    }                                                 
                    if (expertiseSearchedArray[i] == "CustomDevelopment") {
                        expertiseSearchedArray[i] = "Custom Development";                    
                    }                                                  
                    if (expertiseSearchedArray[i] == "PerformanceTuning") {
                        expertiseSearchedArray[i] = "Performance Tuning";                    
                    } 
                  if (expertiseSearchedArray[i] == "SVNtoGitmigration") {
                        expertiseSearchedArray[i] = "SVN to Git migration";                    
                    } 
               
                };                                                        
                  expertiseSearched = expertiseSearchedArray.join(',');

        $('input[name="expertise"]').each(function(i){
          var arrayPosition = jQuery.inArray($(this).attr('value'), expertiseSearched.split(','));

          if (arrayPosition !== -1) {

            $(this).attr('checked','checked');
          }

        });
                   
        $('input[name="products"]').each(function(){

          var arrayPosition = jQuery.inArray($(this).attr('value'), productsSearched.split(','));

          if (arrayPosition !== -1) {
            $(this).attr('checked','checked');
          }

        });


        $("#partnerLocation option").each(function(){
          if (expertLocation == $(this).attr('value')) {
            $(this).attr('selected', 'selected');
          }
          
          
          
        });

        var searchUrl = '/en/wac/resources/partnerSearch.html?' + "expertise=" + expertiseSearched + "&products=" + productsSearched + "&partnerLocation=" + expertLocation + "&isPremier=" + isPremier;
        

        $.get(searchUrl, function(data) {
            $('#resultsAjax').html(data);            
            $('#resultsAjax').pajinate({
              num_page_links_to_display : 10,
              items_per_page : 10,
              item_container_id : '.page-content'  
            });  
            $('#resultsAjax').removeClass('loading');
        });
                         

    $('#searchPartner').submit(function(e){
      e.preventDefault();                                         

      var expertiseCheckedLength = $('input[name="expertise"]:checked').length; 
      var productsCheckedLength = $('input[name="products"]:checked').length;
      var isPremierCheckedLength = $('input[name="ispremier"]:checked').length;

      if (expertiseCheckedLength < 1) {
           var expertiseChecked ='*';             
      } else {
        var expertiseChecked = '';              
      }                             

      if (productsCheckedLength < 1) {
        var productsChecked ='*';
      } else {
        var productsChecked ='';              
      }
                      
      if (isPremierCheckedLength < 1) {
        var isPremierChecked ='*';              
      }  else {
        var isPremierChecked ='true';
      }

      var partnerLocation = $('#partnerLocation option:selected').val();                       

      $('input[name="expertise"]:checked').each(function(i){         

        if ( i != expertiseCheckedLength -1) {
          expertiseChecked = expertiseChecked + $(this).attr('value') + ',';            
        } else {
          expertiseChecked = expertiseChecked + $(this).attr('value');              
        } 

      });

      $('input[name="products"]:checked').each(function(i){         

        if ( i != productsCheckedLength -1) {
          productsChecked = productsChecked + $(this).attr('value') + ',';            
        } else {
          productsChecked = productsChecked + $(this).attr('value');              
        }               

      });

      var searchUrl = '/en/wac/resources/partnerSearch.html?' + "expertise=" + expertiseChecked + "&products=" + productsChecked + "&partnerLocation=" + partnerLocation + "&isPremier=" + isPremierChecked;

            var expertiseCheckedArray = expertiseChecked.split(',');

            for (var i=0; i < expertiseCheckedArray.length; i++) {

                    expertiseCheckedArray[i] = expertiseCheckedArray[i].replace(/\s/g,'');                            
            };                                                        
              expertiseChecked = expertiseCheckedArray.join(',');


      var historyUrl = '?tab=find-an-expert' + "&expertskills=" + expertiseChecked + "&products=" + productsChecked + "&expertlocation=" + partnerLocation + "&isPremier=" + isPremierChecked;
      
      $('#resultsAjax').addClass('loading');
      $('#resultsAjax .featureItem').css('visibility', 'hidden');
      
      setTimeout(function(){                 
        
        $.get(searchUrl, function(data) {
            $('#resultsAjax').html(data);            
          $('#resultsAjax').pajinate({
            num_page_links_to_display : 10,
            items_per_page : 10,
            item_container_id : '.page-content'  
          }); 
          $('#resultsAjax').removeClass('loading');
          $('#resultsAjax .featureItem').css('visibility', 'visible');
        });
        
      }, 500);

                    History.pushState(null, "Experts", historyUrl);


    });
  }
  expertSearchResults();
}

var gPlus = function(){                
  $(document).ajaxStop(function(){
    gapi.plusone.go('googlePlusOnDemand');
      gapi.plusone.go('googlePlusDownload');    
  });

}

var jobPostingNumbers = function(){
    var numJobsArray = [];
    $(".tabsContent").each(function(i){
        var job = $(this).find('li');
        z = 1;
        for (var x=0; x < job.length; x++) {
            if (job.eq(x).css('display') == 'block') {
                numJobsArray[i] = z++;
            }
        };
    });                                                

     $(".numOfJobs").each(function(i){
     if (numJobsArray[i] == undefined) {
      $(this).html(" (0) ");
    } else {
      $(this).html(" (" +numJobsArray[i] + ") ");
    }
    
     });
     
}

var filterFeaturedRecent = function(settings){
    if (!settings) {
        var settings = {
            listClassName : '.test',
      limitToTenArticles : true
        };
    }                                           

  if (settings.limitToTenArticles == true) {
    $(settings.listClassName +' li.featured').each(function(i){
      if (i>9) {           
        $(this).hide();
      }
    });
  }


    $('.filterByFeatured').click(function(e){
        e.preventDefault();
        $(settings.listClassName +' li').hide();    
        $(settings.listClassName +' li.featured').show();
        $('.featuredFilter').find('a').removeClass('selected');
        $(this).addClass('selected');         
    if (settings.limitToTenArticles == true) {
      $(settings.listClassName +' li.featured').each(function(i){
        if (i>9) {           
          $(this).hide();
        }
      });
    }
    });                    
    $('.filterByRecent').click(function(e){

        e.preventDefault();          
        $('.featuredFilter').find('a').removeClass('selected');
        $(this).addClass('selected');            
        $(settings.listClassName +' li').show();  
    if (settings.limitToTenArticles == true) {          
      $(settings.listClassName +' li').each(function(i){
        if (i>9) {           
          $(this).hide();
        }
      });                                       
    }
    });
}

// Product Overview Tooltip

var toolTipInit = function() {
  
  $('.tooltip-label').hover(function() {
    var toolTipWidth = $(this).outerWidth() + 22;
    if($(this).parent().parent().hasClass('imageRight')) {
      $(this).next('.tooltip').css({'left':toolTipWidth});
    }
    $(this).next('.tooltip').stop(true, true).fadeIn(200);
  }, function() {
    $(this).next('.tooltip').stop(true, true).fadeOut(200);
  });
}


var merlinTrack = function(category, action, label){
    $.ajax({ url: "/merlinTrack?_mcategory=" + category + "&_maction=" + action + "&_mlabel=" + label });
}