const qlist = new Vue({
    el: '#qlist',
    data: {
      // search condition
      questionId : '',      // 検索条件として設定しているquestionId
      ownerGroupId: 1,      // 検索条件として設定しているonwerGroupId（今のところ固定)
      selectedLesson : 0,   // 検索条件として選択しているLesson。暫定初期値として 「"0":全て」を選択状態とする

      // selected question list
      lessons : [],         // 検索条件として表示されるLessonの一覧(ownerGroupIdに紐づくLessonが設定される)

      // question list
      questions : []        // 検索結果として表示される問題の一覧
    },
    methods: {
      init: () => {
        const payload = {
          'ownerGroupID' : Number(qlist.ownerGroupId),
          'selectedLesson' : Number(qlist.selectedLesson)
        }
        // 出題用の問題を取得し、qlistオブジェクトに格納する。
        axios.post('api/qlist/init', payload)
        .then((response) => {
          //alert(JSON.stringify(response));
          // 検索対象選択用のLessonリストを初期化
          qlist.lessons[0] = {
            "lessonId" : 0,
            "lessonName" : "全て"
          }
          response.data.lessons.forEach(function(v,i) {
            qlist.lessons[i+1] = { 
              "lessonId" : v.lessonId,
              "lessonName" : v.lessonName
            }
            //alert("lessonId:" + qlist.lessons[i+1].lessonId + ",lessonName:" + qlist.lessons[i+1].lessonName);
          });
          // ↑でqlist.lessonsの初期化前にselectedLessonを設定する画面の初期表示時に選択肢が正しく設定されない。
          // qlist.lessons→qlist.selectedLessonの順序で設定する必要がある。
          // 選択した検索条件(lesson)を設定
          qlist.selectedLesson = String(response.data.searchCondition.selectedLesson);
        })
        .catch((err) => {
          alert(err)
        })
      },
      search: () => {
        const payload = {
          'ownerGroupID' : Number(qlist.ownerGroupId),
          'selectedLesson' : Number(qlist.selectedLesson)
        }
        // 出題用の問題を取得し、qlistオブジェクトに格納する。
        axios.post('api/qlist/list', payload)
        .then((response) => {
          //alert(JSON.stringify(response));
          qlist.questions = response.data.questions;
        })
        .catch((err) => {
          alert(err)
        })
      },
      // 画面ロード時にinitメソッドを実行
      window:onload = function() {
        qlist.init();
      },
    }
  })