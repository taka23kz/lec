const quiz = new Vue({
    el: '#quiz',
    data: {
      ownerGroupID: 1,
      question : {
        questionId : '',
        question : '',    // 画面に表示する問題文
        answerType : '',
        choiceNum : '',
        ownerGroupId : '',
        lessonId : '',
      },
      choiceId : '',
      choiceIds : [],
      choices : []      // 画面に表示する選択肢のリスト
    },
    methods: {
      // 出題用のquiz情報を初期化する。
      quiz: () => {
        const payload = {
          'ownerGroupID' : quiz.ownerGroupID
        }
        // 出題用の問題を取得し、quizオブジェクトに格納する。
        axios.post('api/question/quiz', payload)
        .then((response) => {
          //alert(JSON.stringify(response));
          //alert(response.data.Question.question);
          quiz.question.questionId = response.data.Question.questionId;
          quiz.question.question = response.data.Question.question;
          quiz.question.answerType = response.data.Question.answerType;

          response.data.Choices.forEach(function(v,i) {
            quiz.choices[i] = {
              "choiceId" : v.choiceId,
              "questionId" : v.questionId,
              "choiceLabel" : v.choiceLabel,
              "correct" : v.correct
            }
          });
        })
        .catch((err) => {
          alert(err.response.data.error)
        })
      },
      // 回答ボタン押下時のアクション
      answer: () => {
        var payload;
        // 回答の形式によってサービスに渡すパラメータが異なるためanswerTypeの値でリクエストを編集
        if ( quiz.question.answerType == "01" ) {
          // ラジオボタン形式
          payload = {
            'questionID' : quiz.question.questionId,
            'answerType' : quiz.question.answerType,
            'choice' : quiz.choiceId
          }
        } else if( quiz.question.answerType == "02" ) {
          // チェックボックス形式
          payload = {
            'questionID' : quiz.question.questionId,
            'answerType' : quiz.question.answerType,
            'choiceIds' : quiz.choiceIds
          }
        }
        axios.post('api/question/answer', payload)
        .then((response) => {
          if ( response.data.correct ) {
            // 正解時のアクション
            alert('正解！');
          } else {
            // 不正解時のアクション
            alert('不正解!');
          }
          //alert(response.data.answer);
        }).catch((err) => {
          alert(err.response.data.error)
        })
      },
      // 画面ロード時にquizメソッドを実行
      window:onload = function() {
        quiz.quiz();
      },
    }
  })