const signup = new Vue({
    el: '#signup',
    data: {
      // 登録内容
      mailAddress : null, 
      userName : null,      // 登録するユーザ名
      passwd : null,      // 登録するユーザのパスワード
      passwd2 : null      // 登録するユーザのパスワード(確認用)
    },
    methods: {
      signup: () => {
        const payload = {
          'mailAddress' : signup.mailAddress,
          'userName' : signup.userName,
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