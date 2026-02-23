package repository

import (
	"sort"
	"sync"
	"time"

	"asenare/backend/internal/domain"
)

type InMemoryGlossaryRepository struct {
	mu    sync.RWMutex
	terms map[string]domain.GlossaryTerm
}

func NewInMemoryGlossaryRepository() *InMemoryGlossaryRepository {
	now := time.Now().UTC()
	seed := []domain.GlossaryTerm{
		// FP3級関連
		{ID: "g01", Term: "ライフプラン", Reading: "らいふぷらん", Definition: "人生の各段階における目標を設定し、それを実現するための長期的な生活設計。", CreatedAt: now, UpdatedAt: now},
		{ID: "g02", Term: "キャッシュフロー", Reading: "きゃっしゅふろー", Definition: "一定期間における現金の流入と流出の状況。収入から支出を引いた差額のこと。", CreatedAt: now, UpdatedAt: now},
		{ID: "g03", Term: "純資産", Reading: "じゅんしさん", Definition: "資産合計から負債合計を差し引いた正味の資産。自己資本とも呼ばれる。", CreatedAt: now, UpdatedAt: now},
		{ID: "g04", Term: "学資保険", Reading: "がくしほけん", Definition: "子どもの教育資金を積み立てるための保険。満期時に祝い金・満期金を受け取れる貯蓄性のある保険。", CreatedAt: now, UpdatedAt: now},
		{ID: "g05", Term: "住宅ローン控除", Reading: "じゅうたくろーんこうじょ", Definition: "住宅を購入した際に住宅ローン残高の一定割合を所得税から控除できる制度。最長13年間適用される。", CreatedAt: now, UpdatedAt: now},
		{ID: "g06", Term: "生命保険", Reading: "せいめいほけん", Definition: "被保険者が死亡または高度障害になった際に保険金が支払われる保険。定期保険・終身保険・養老保険などがある。", CreatedAt: now, UpdatedAt: now},
		{ID: "g07", Term: "定期保険", Reading: "ていきほけん", Definition: "一定期間のみ死亡保障が続く保険。掛け捨て型で保険料が安いが、満期時には保険金が支払われない。", CreatedAt: now, UpdatedAt: now},
		{ID: "g08", Term: "終身保険", Reading: "しゅうしんほけん", Definition: "一生涯にわたって死亡保障が続く保険。貯蓄性があり解約返戻金がある。", CreatedAt: now, UpdatedAt: now},
		{ID: "g09", Term: "国民年金", Reading: "こくみんねんきん", Definition: "日本に住む20歳以上60歳未満の全ての人が加入する公的年金制度の基礎部分。老齢・障害・遺族の3つの給付がある。", CreatedAt: now, UpdatedAt: now},
		{ID: "g10", Term: "厚生年金", Reading: "こうせいねんきん", Definition: "会社員・公務員が国民年金に上乗せして加入する公的年金。報酬比例部分があり、収入に応じた給付を受けられる。", CreatedAt: now, UpdatedAt: now},
		{ID: "g11", Term: "ペイオフ", Reading: "ぺいおふ", Definition: "金融機関が破綻した際、預金者1人あたり元本1,000万円とその利息までを預金保険機構が保護する制度。", CreatedAt: now, UpdatedAt: now},
		{ID: "g12", Term: "複利", Reading: "ふくり", Definition: "利息を元本に加算し、その合計に対してさらに利息を計算する方式。長期運用で元本に対する利息が雪だるま式に増える。", CreatedAt: now, UpdatedAt: now},
		{ID: "g13", Term: "ドルコスト平均法", Reading: "どるこすとへいきんほう", Definition: "一定期間ごとに一定金額で投資商品を購入する方法。価格が高い時は少なく、安い時は多く買えるため平均購入コストを抑えられる。", CreatedAt: now, UpdatedAt: now},
		{ID: "g14", Term: "NISA", Reading: "にーさ", Definition: "少額投資非課税制度。NISA口座内で購入した株式・投資信託等の売却益や配当金が非課税になる制度。", CreatedAt: now, UpdatedAt: now},
		{ID: "g15", Term: "iDeCo", Reading: "いでこ", Definition: "個人型確定拠出年金。掛け金が全額所得控除・運用益が非課税・受取時も控除ありの3つの税制優遇がある老後資産形成制度。", CreatedAt: now, UpdatedAt: now},
		{ID: "g16", Term: "所得控除", Reading: "しょとくこうじょ", Definition: "課税所得を計算する際に所得金額から差し引ける金額。基礎控除・配偶者控除・医療費控除などがある。", CreatedAt: now, UpdatedAt: now},
		{ID: "g17", Term: "ふるさと納税", Reading: "ふるさとのうぜい", Definition: "自分の選んだ自治体に寄附することで、実質2,000円の自己負担で返礼品がもらえ、寄附額から2,000円を引いた金額が所得税・住民税から控除される制度。", CreatedAt: now, UpdatedAt: now},
		{ID: "g18", Term: "相続税", Reading: "そうぞくぜい", Definition: "被相続人（亡くなった人）から財産を相続した場合にかかる税金。基礎控除は3,000万円＋600万円×法定相続人の数。", CreatedAt: now, UpdatedAt: now},
		{ID: "g19", Term: "贈与税", Reading: "ぞうよぜい", Definition: "個人から財産をもらった場合にかかる税金。年間110万円の基礎控除がある暦年課税と、2,500万円まで非課税の相続時精算課税がある。", CreatedAt: now, UpdatedAt: now},
		// 簿記3級関連
		{ID: "g20", Term: "仕訳", Reading: "しわけ", Definition: "簿記における取引を借方（左）と貸方（右）に分類して記録する作業。借方と貸方の金額は常に一致する。", CreatedAt: now, UpdatedAt: now},
		{ID: "g21", Term: "勘定科目", Reading: "かんじょうかもく", Definition: "簿記で取引内容を分類するための名称。資産・負債・純資産・収益・費用の5つに分類される。", CreatedAt: now, UpdatedAt: now},
		{ID: "g22", Term: "売掛金", Reading: "うりかけきん", Definition: "商品を掛けで販売したときに生じる代金の受取権利。後日回収する資産。", CreatedAt: now, UpdatedAt: now},
		{ID: "g23", Term: "買掛金", Reading: "かいかけきん", Definition: "商品を掛けで仕入れたときに生じる代金の支払義務。後日支払う負債。", CreatedAt: now, UpdatedAt: now},
		{ID: "g24", Term: "減価償却", Reading: "げんかしょうきゃく", Definition: "固定資産の取得原価を耐用年数にわたって費用として配分する会計処理。土地は価値が減らないため対象外。", CreatedAt: now, UpdatedAt: now},
		{ID: "g25", Term: "損益計算書", Reading: "そんえきけいさんしょ", Definition: "一定期間の収益と費用を対比して当期純利益または純損失を示す財務諸表。P/Lとも呼ばれる。", CreatedAt: now, UpdatedAt: now},
		{ID: "g26", Term: "貸借対照表", Reading: "たいしゃくたいしょうひょう", Definition: "一定時点における企業の資産・負債・純資産を示す財務諸表。B/Sとも呼ばれる。資産合計＝負債合計＋純資産合計が成立する。", CreatedAt: now, UpdatedAt: now},
		// 資産運用関連
		{ID: "g27", Term: "配当金", Reading: "はいとうきん", Definition: "企業が利益の一部を株主に分配するお金。配当利回りは1株当たり配当額を株価で割った値。", CreatedAt: now, UpdatedAt: now},
		{ID: "g28", Term: "投資信託", Reading: "とうししんたく", Definition: "多くの投資家から集めた資金をまとめて専門家が分散投資する金融商品。ファンドとも呼ばれる。", CreatedAt: now, UpdatedAt: now},
		{ID: "g29", Term: "ETF", Reading: "いーてぃーえふ", Definition: "Exchange Traded Fund（上場投資信託）の略。株式市場に上場し、株式と同様にリアルタイムで売買できる投資信託。", CreatedAt: now, UpdatedAt: now},
		{ID: "g30", Term: "ポートフォリオ", Reading: "ぽーとふぉりお", Definition: "保有している資産の組み合わせ。複数の資産に分散投資することでリスクを軽減できる。", CreatedAt: now, UpdatedAt: now},
		{ID: "g31", Term: "PER", Reading: "ぴーいーあーる", Definition: "株価収益率（Price Earnings Ratio）。株価を1株当たり純利益で割った値。低いほど株価が割安とされる。", CreatedAt: now, UpdatedAt: now},
		{ID: "g32", Term: "インデックスファンド", Reading: "いんでっくすふぁんど", Definition: "日経平均やS&P500などの市場指数に連動することを目指す投資信託。コストが低く長期投資に適している。", CreatedAt: now, UpdatedAt: now},
		{ID: "g33", Term: "信用リスク", Reading: "しんようりすく", Definition: "債券などの発行体が利息の支払いや元本の返済ができなくなるリスク（デフォルトリスク）。", CreatedAt: now, UpdatedAt: now},
	}
	m := make(map[string]domain.GlossaryTerm, len(seed))
	for _, t := range seed {
		m[t.ID] = t
	}
	return &InMemoryGlossaryRepository{terms: m}
}

func (r *InMemoryGlossaryRepository) List() ([]domain.GlossaryTerm, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	list := make([]domain.GlossaryTerm, 0, len(r.terms))
	for _, t := range r.terms {
		list = append(list, t)
	}
	sort.Slice(list, func(i, j int) bool {
		return list[i].Reading < list[j].Reading
	})
	return list, nil
}

func (r *InMemoryGlossaryRepository) FindByID(id string) (domain.GlossaryTerm, bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	t, ok := r.terms[id]
	return t, ok, nil
}
