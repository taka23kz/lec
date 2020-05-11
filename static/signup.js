const signup = new Vue({
    el: '#signup',
    data: {
      // 登録内容
      userId : null,
      userName : null,      // 登録するユーザ名
      mailAddress : null, 
      passwd : null,      // 登録するユーザのパスワード
      passwd2 : null      // 登録するユーザのパスワード(確認用)
    },
    methods: {
      signup: () => {
        const payload = {
          'userId'  : signup.userId,
          'userName' : signup.userName,
          'mailAddress' : signup.mailAddress,
          'passwd' : signup.passwd,
          'passwd2' : signup.passwd2
        }
        axios.post('api/signup/signup', payload)
        .then((response) => {
          alert(JSON.stringify(response));
        })
        .catch((err) => {
          alert(err)
        })
      },
      // 画面ロード時にメソッドを実行
      window:onload = function() {
      },
    }
  })