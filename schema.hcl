
schema "coastline" {
  comment = "A schema comment"
}

table "project" {
  schema = schema.example
  column "id" {
    id {
      generated = ALWAYS
    }
    type = int
    null = false
  }
  column "name" {
    type = text
    null = false
  }

  primary_key {
    columns = [column.id]
  }
}


