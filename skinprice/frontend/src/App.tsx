import React from "react";
import "./styles.css";

type StatItemType = {
  icon: string;
  price: string;
  pricePercent: string;
  quantity: string;
  quantityPercent: string;
};

type CardData = {
  image: string;
  title: string;
  stats: StatItemType[];
};

const StatItem: React.FC<StatItemType> = ({
  icon,
  price,
  pricePercent,
  quantity,
  quantityPercent,
}) => {
  return (
    <div className="stat-item">
      <img src={icon} alt="store" className="store-icon" />

      <div className="stat-content">
        <div className="row">
          <span className="text">{price}</span>
          <span className="positive">{pricePercent}</span>

          <span className="feature">
            ликвидность <span className="positive">+12%</span>
          </span>
        </div>

        <div className="row">
          <span className="text">{quantity}</span>
          <span className="negative">{quantityPercent}</span>

          <span className="feature">
            спрос <span className="negative">-3%</span>
          </span>
        </div>
      </div>
    </div>
  );
};

type CardProps = {
  data: CardData;
};

const Card: React.FC<CardProps> = ({ data }) => {
  return (
    <div className="card">
      <div className="image-wrapper">
        <img src={data.image} alt="weapon" className="card-image" />
      </div>

      <div className="card-body">
        <h2 className="title">{data.title}</h2>

        <div className="stats-list">
          {data.stats.map((stat, index) => (
            <StatItem key={index} {...stat} />
          ))}
        </div>
      </div>
    </div>
  );
};

/* ===== ДАННЫЕ ===== */

const stores = [
  "https://cdn-icons-png.flaticon.com/512/5968/5968705.png",
  "https://cdn-icons-png.flaticon.com/512/5968/5968672.png",
  "https://cdn-icons-png.flaticon.com/512/5968/5968866.png",
];

const cards: CardData[] = [
  {
    image:
      "https://images.unsplash.com/photo-1542751371-adc38448a05e?auto=format&fit=crop&w=900&q=80",
    title: "MAC-10 | BEBRA",
    stats: [
      {
        icon: stores[0],
        price: "1000 ₽",
        pricePercent: "+10%",
        quantity: "5693 шт.",
        quantityPercent: "-5%",
      },
      {
        icon: stores[1],
        price: "1430 ₽",
        pricePercent: "+24%",
        quantity: "2512 шт.",
        quantityPercent: "-1%",
      },
      {
        icon: stores[2],
        price: "840 ₽",
        pricePercent: "+3%",
        quantity: "9231 шт.",
        quantityPercent: "-11%",
      },
    ],
  },

  {
    image:
      "https://images.unsplash.com/photo-1511512578047-dfb367046420?auto=format&fit=crop&w=900&q=80",
    title: "AK-47 | FIRE",
    stats: [
      {
        icon: stores[0],
        price: "2100 ₽",
        pricePercent: "+6%",
        quantity: "1200 шт.",
        quantityPercent: "-2%",
      },
      {
        icon: stores[1],
        price: "1900 ₽",
        pricePercent: "+2%",
        quantity: "980 шт.",
        quantityPercent: "-4%",
      },
      {
        icon: stores[2],
        price: "2500 ₽",
        pricePercent: "+12%",
        quantity: "540 шт.",
        quantityPercent: "-8%",
      },
    ],
  },

  {
    image:
      "https://images.unsplash.com/photo-1511512578047-dfb367046420?auto=format&fit=crop&w=900&q=80",
    title: "AK-47 | FIRE",
    stats: [
      {
        icon: stores[0],
        price: "2100 ₽",
        pricePercent: "+6%",
        quantity: "1200 шт.",
        quantityPercent: "-2%",
      },
      {
        icon: stores[1],
        price: "1900 ₽",
        pricePercent: "+2%",
        quantity: "980 шт.",
        quantityPercent: "-4%",
      },
      {
        icon: stores[2],
        price: "2500 ₽",
        pricePercent: "+12%",
        quantity: "540 шт.",
        quantityPercent: "-8%",
      },
    ],
  },

  {
    image:
      "https://images.unsplash.com/photo-1511512578047-dfb367046420?auto=format&fit=crop&w=900&q=80",
    title: "AK-47 | FIRE",
    stats: [
      {
        icon: stores[0],
        price: "2100 ₽",
        pricePercent: "+6%",
        quantity: "1200 шт.",
        quantityPercent: "-2%",
      },
      {
        icon: stores[1],
        price: "1900 ₽",
        pricePercent: "+2%",
        quantity: "980 шт.",
        quantityPercent: "-4%",
      },
      {
        icon: stores[2],
        price: "2500 ₽",
        pricePercent: "+12%",
        quantity: "540 шт.",
        quantityPercent: "-8%",
      },
    ],
  },

  {
    image:
      "https://images.unsplash.com/photo-1511512578047-dfb367046420?auto=format&fit=crop&w=900&q=80",
    title: "AK-47 | FIRE",
    stats: [
      {
        icon: stores[0],
        price: "2100 ₽",
        pricePercent: "+6%",
        quantity: "1200 шт.",
        quantityPercent: "-2%",
      },
      {
        icon: stores[1],
        price: "1900 ₽",
        pricePercent: "+2%",
        quantity: "980 шт.",
        quantityPercent: "-4%",
      },
      {
        icon: stores[2],
        price: "2500 ₽",
        pricePercent: "+12%",
        quantity: "540 шт.",
        quantityPercent: "-8%",
      },
    ],
  },

  {
    image:
      "https://images.unsplash.com/photo-1511512578047-dfb367046420?auto=format&fit=crop&w=900&q=80",
    title: "AK-47 | FIRE",
    stats: [
      {
        icon: stores[0],
        price: "2100 ₽",
        pricePercent: "+6%",
        quantity: "1200 шт.",
        quantityPercent: "-2%",
      },
      {
        icon: stores[1],
        price: "1900 ₽",
        pricePercent: "+2%",
        quantity: "980 шт.",
        quantityPercent: "-4%",
      },
      {
        icon: stores[2],
        price: "2500 ₽",
        pricePercent: "+12%",
        quantity: "540 шт.",
        quantityPercent: "-8%",
      },
    ],
  },

  {
    image:
      "https://images.unsplash.com/photo-1511512578047-dfb367046420?auto=format&fit=crop&w=900&q=80",
    title: "AK-47 | FIRE",
    stats: [
      {
        icon: stores[0],
        price: "2100 ₽",
        pricePercent: "+6%",
        quantity: "1200 шт.",
        quantityPercent: "-2%",
      },
      {
        icon: stores[1],
        price: "1900 ₽",
        pricePercent: "+2%",
        quantity: "980 шт.",
        quantityPercent: "-4%",
      },
      {
        icon: stores[2],
        price: "2500 ₽",
        pricePercent: "+12%",
        quantity: "540 шт.",
        quantityPercent: "-8%",
      },
    ],
  },

  {
    image:
      "https://images.unsplash.com/photo-1511512578047-dfb367046420?auto=format&fit=crop&w=900&q=80",
    title: "AK-47 | FIRE",
    stats: [
      {
        icon: stores[0],
        price: "2100 ₽",
        pricePercent: "+6%",
        quantity: "1200 шт.",
        quantityPercent: "-2%",
      },
      {
        icon: stores[1],
        price: "1900 ₽",
        pricePercent: "+2%",
        quantity: "980 шт.",
        quantityPercent: "-4%",
      },
      {
        icon: stores[2],
        price: "2500 ₽",
        pricePercent: "+12%",
        quantity: "540 шт.",
        quantityPercent: "-8%",
      },
    ],
  },

  {
    image:
      "https://images.unsplash.com/photo-1511512578047-dfb367046420?auto=format&fit=crop&w=900&q=80",
    title: "AK-47 | FIRE",
    stats: [
      {
        icon: stores[0],
        price: "2100 ₽",
        pricePercent: "+6%",
        quantity: "1200 шт.",
        quantityPercent: "-2%",
      },
      {
        icon: stores[1],
        price: "1900 ₽",
        pricePercent: "+2%",
        quantity: "980 шт.",
        quantityPercent: "-4%",
      },
      {
        icon: stores[2],
        price: "2500 ₽",
        pricePercent: "+12%",
        quantity: "540 шт.",
        quantityPercent: "-8%",
      },
    ],
  },

  {
    image:
      "https://images.unsplash.com/photo-1511512578047-dfb367046420?auto=format&fit=crop&w=900&q=80",
    title: "AK-47 | FIRE",
    stats: [
      {
        icon: stores[0],
        price: "2100 ₽",
        pricePercent: "+6%",
        quantity: "1200 шт.",
        quantityPercent: "-2%",
      },
      {
        icon: stores[1],
        price: "1900 ₽",
        pricePercent: "+2%",
        quantity: "980 шт.",
        quantityPercent: "-4%",
      },
      {
        icon: stores[2],
        price: "2500 ₽",
        pricePercent: "+12%",
        quantity: "540 шт.",
        quantityPercent: "-8%",
      },
    ],
  },
];

const App: React.FC = () => {
  return (
    <div className="app">
      <div className="container">
        {cards.map((card, index) => (
          <Card key={index} data={card} />
        ))}
      </div>
    </div>
  );
};

export default App;